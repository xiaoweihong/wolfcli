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

func GetVolumeInfo(day string) (volumeList []*VolumeServerMap) {
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		fmt.Println(err)
		return
	}
	var mc *wdclient.MasterClient
	var resp *master_pb.VolumeListResponse
	mc = wdclient.NewMasterClient(grpc.WithInsecure(), "client", "", 0,
		strings.Split(HostAndPort(), ","))
	go mc.KeepConnectedToMaster()
	mc.WaitUntilConnected()
	err = mc.WithClient(func(client master_pb.SeaweedClient) (err error) {
		resp, err = client.VolumeList(context.Background(), &master_pb.VolumeListRequest{})
		return nil
	})
	if err != nil {
		fmt.Println(err)
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
			if TimeNowZero().Sub(unix) > time.Duration(time.Hour*24*time.Duration(dayInt)) {
				fmt.Printf("volumeId:%v\tsize:%v\ttime:%v ttl:%v\n",
					volume.GetId(),
					humanize.Bytes(volume.GetSize()),
					unix,
					volume.GetTtl())
				volumeServerMap.Volume = append(volumeServerMap.Volume, volume.GetId())
			}
		}
		volumeList = append(volumeList, volumeServerMap)
	}
	return
}

func DeleteVolumeById(id string) {
	volumeId, err := needle.NewVolumeId(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	listVolume := GetVolumeInfo(global.Day)
	tid, _ := strconv.Atoi(id)
	for _, server := range listVolume {
		if IsContain(server.Volume, uint32(tid)) {
			err := deleteVolume(grpc.WithInsecure(), volumeId, server.VolumeServer)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%v卷删除成功\n", tid)
			}
		} else {
			fmt.Printf("%v卷不存在集群中，请检查输入的卷id\n", tid)
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
