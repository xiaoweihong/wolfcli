package controller

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb/master_pb"
	"github.com/chrislusf/seaweedfs/weed/pb/volume_server_pb"
	"github.com/chrislusf/seaweedfs/weed/wdclient"
	"github.com/dustin/go-humanize"
	"google.golang.org/grpc"
	"sort"
	"strconv"
	"strings"
	"time"
	"wolfcli/global"
)

type VolumeServerMap struct {
	VolumeServer string
	Volume       []uint32
}

var (
	mc   *wdclient.MasterClient
	resp *master_pb.VolumeListResponse
)

// GetVolumeInfo 获取集群内图片信息
func GetVolumeInfo(day string) (volumeList []*VolumeServerMap) {
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := CheckServer(); err != nil {
		fmt.Println("seaweedfs连不上,请检查", err)
		return
	}
	mc = wdclient.NewMasterClient(
		grpc.WithInsecure(),
		"xiaowei",
		"",
		0,
		"",
		strings.Split(HostAndPort(), ","))
	go mc.KeepConnectedToMaster()
	mc.WaitUntilConnected()
	err = mc.WithClient(func(client master_pb.SeaweedClient) (err error) {
		resp, err = client.VolumeList(context.Background(), &master_pb.VolumeListRequest{})
		return err
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(resp.GetTopologyInfo().GetDataCenterInfos()) == 0 {
		fmt.Println("volume信息不存在，请检查volume与master之间是否正常")
		return
	}
	for _, node := range resp.GetTopologyInfo().GetDataCenterInfos()[0].GetRackInfos()[0].GetDataNodeInfos() {
		volumeServerMap := &VolumeServerMap{}
		volumeServerMap.VolumeServer = node.GetId()
		fmt.Println(strings.Repeat("-~-", 20))
		fmt.Printf("volume server:%v\n", node.GetId())
		fmt.Println(strings.Repeat("-~-", 20))
		volumeInfos := node.GetVolumeInfos()
		sort.Slice(volumeInfos, func(i, j int) bool {
			return volumeInfos[i].ModifiedAtSecond > volumeInfos[j].ModifiedAtSecond
		})

		for _, volume := range volumeInfos {
			unix := time.Unix(volume.GetModifiedAtSecond(), 0)
			if TimeNowZero().Sub(unix) > time.Duration(time.Hour*24*time.Duration(dayInt-1)) {
				fmt.Printf("volumeId:%v\tsize:%v\ttime:%v ttl:%v\tcount:%v\n",
					volume.GetId(),
					humanize.Bytes(volume.GetSize()),
					unix,
					volume.GetTtl(),
					volume.GetFileCount())
				volumeServerMap.Volume = append(volumeServerMap.Volume, volume.GetId())
			}
		}
		volumeList = append(volumeList, volumeServerMap)
	}
	return
}

// DeleteVolumeById 根据volumeId删除volume
func DeleteVolumeById(id string) {

	listVolume := GetVolumeInfo(global.Day)
	tid, _ := strconv.Atoi(id)
	for _, server := range listVolume {
		if IsContain(server.Volume, uint32(tid)) {
			// [{10.244.0.117:9300 192.168.2.174:9300 }]
			//        Url                PublicUrl
			// 找到volume对应的volume server
			locations, _ := mc.GetVidLocations(id)
			//err := DeleteVolumeByVolumeServer(id, locations[0].PublicUrl)
			err := operation.WithVolumeServerClient(locations[0].PublicUrl, grpc.WithInsecure(), func(volumeServerClient volume_server_pb.VolumeServerClient) error {
				_, deleteErr := volumeServerClient.VolumeDelete(context.Background(), &volume_server_pb.VolumeDeleteRequest{
					VolumeId: uint32(tid),
				})
				return deleteErr
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%v卷删除成功\n", tid)
			}
		} else {
			fmt.Printf("%v卷不存在集群中，请检查输入的卷id和天数\n", tid)
		}
	}
}

// DeleteVolumeByTime 根据天数删除volume
func DeleteVolumeByTime() {
	listVolume := GetVolumeInfo(global.Day)
	for _, volume := range listVolume {
		if len(listVolume[0].Volume) == 0 {
			fmt.Printf("volume server->%v\t%v天之前的volume已经删除，暂无可删除的volume\n", volume.VolumeServer, global.Day)
		}
		for _, id := range volume.Volume {
			locations, _ := mc.GetVidLocations(strconv.Itoa(int(id)))
			//err := DeleteVolumeByVolumeServer(strconv.Itoa(int(id)), locations[0].PublicUrl)
			err := operation.WithVolumeServerClient(locations[0].PublicUrl, grpc.WithInsecure(), func(volumeServerClient volume_server_pb.VolumeServerClient) error {
				_, deleteErr := volumeServerClient.VolumeDelete(context.Background(), &volume_server_pb.VolumeDeleteRequest{
					VolumeId: uint32(id),
				})
				return deleteErr
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%v卷删除成功\n", id)
			}

		}
	}

}

// HostAndPort 返回master server 格式的地址.
// 127.0.0.1:9333
func HostAndPort() string {
	return global.PicIP.String() + ":" + global.PicPort
}

// TimeNowZero 获得当前时间的0点.
// 2021-02-18 15:55:22 +0800 CST --> 2021-02-18 00:00:00 +0800 CST
func TimeNowZero() time.Time {
	parse, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	return parse
}

// IsContain 查询数字item是否存在于切片items里.
func IsContain(items []uint32, item uint32) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// CheckServer 检查seaweedfsserver是否可联通.
func CheckServer() error {
	get := httplib.Get(fmt.Sprintf("http://%s", HostAndPort()))
	get.SetTimeout(time.Second*3, time.Second)
	_, err := get.Bytes()
	if err != nil {
		return err
	}
	return nil
}
