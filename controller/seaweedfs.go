package controller

import (
	"context"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb/master_pb"
	"github.com/chrislusf/seaweedfs/weed/pb/volume_server_pb"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
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

func GetVolumeInfo(day string) (volumeList []*VolumeServerMap) {
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		fmt.Println(err)
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

func DeleteVolumeById(id string) {
	//volumeId, err := needle.NewVolumeId(id)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	listVolume := GetVolumeInfo(global.Day)
	tid, _ := strconv.Atoi(id)
	for _, server := range listVolume {
		if IsContain(server.Volume, uint32(tid)) {
			//vServer := server.VolumeServer
			//err := deleteVolume(grpc.WithInsecure(), volumeId, vServer)
			locations, _ := mc.GetVidLocations(id)
			err := DeleteVolume(id, locations[0].PublicUrl)
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

func DeleteVolume(id string, volumeServer string) error {
	volumeId, err := needle.NewVolumeId(id)
	if err != nil {
		return err
	}
	err = deleteVolume(grpc.WithInsecure(), volumeId, volumeServer)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func DeleteVolumeByTime() {
	listVolume := GetVolumeInfo(global.Day)
	for _, volume := range listVolume {
		if len(listVolume[0].Volume) == 0 {
			fmt.Printf("volume server->%v\t%v天之前的volume已经删除，暂无可删除的volume\n", volume.VolumeServer, global.Day)
		}
		for _, id := range volume.Volume {
			locations, _ := mc.GetVidLocations(strconv.Itoa(int(id)))
			err := DeleteVolume(strconv.Itoa(int(id)), locations[0].PublicUrl)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%v卷删除成功\n", id)
			}

		}
	}

}

func deleteVolume(grpcDialOption grpc.DialOption, volumeId needle.VolumeId, sourceVolumeServer string) (err error) {
	return operation.WithVolumeServerClient(sourceVolumeServer, grpcDialOption, func(volumeServerClient volume_server_pb.VolumeServerClient) error {
		_, deleteErr := volumeServerClient.VolumeDelete(context.Background(), &volume_server_pb.VolumeDeleteRequest{
			VolumeId: uint32(volumeId),
		})
		return deleteErr
	})
}

func HostAndPort() string {
	return global.PicIP.String() + ":" + global.PicPort
}

func TimeNowZero() time.Time {
	parse, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	return parse
}

func IsContain(items []uint32, item uint32) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
