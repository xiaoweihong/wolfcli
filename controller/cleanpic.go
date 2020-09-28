package controller

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"wolfcli/global"
	"wolfcli/model"
	"wolfcli/utils"
)

var (
	picCount int64
)

func GetFaceEntites(startTime int64, endTime int64) ([]*model.FaceTable, error) {
	var faces []*model.FaceTable
	err := global.Db.Table(model.FaceTableName).
		Where("ts >? ", startTime).
		Where("ts<?", endTime).
		Find(&faces).Error
	if err != nil {
		zap.L().Sugar().Error("获取数据错误")
		return nil, err
	}
	return faces, err
}

func DeleteFaces(startTime int64, endTime int64) {
	totalCount := 0
	tmpStartTime := endTime - 3600000
	if tmpStartTime < 0 {
		tmpStartTime = startTime
	}
	tmpEndTime := endTime
	k := 0
	deleteThreadChan := make(chan int, global.ParamNum)
	for i := 0; i < global.ParamNum; i++ {
		deleteThreadChan <- i
	}
	client := &http.Client{
		Timeout: global.HttpTimeOut * time.Second,
	}
	for {
		zap.L().Sugar().Info("开始第", k, "次循环", "要删除起始时间：", tmpStartTime, "要删除结束时间：", tmpEndTime)
		wg := sync.WaitGroup{}
		if tmpStartTime >= tmpEndTime || tmpEndTime < startTime {
			zap.L().Sugar().Info("开始压缩weed")
			err := weedShrink(client)
			if err != nil {
				fmt.Println("压缩weed失败：", err.Error())
			}
			zap.L().Sugar().Info("清理完成！共清理数据:", totalCount)
			return
		}
		faceEntities, err := GetFaceEntites(tmpStartTime, tmpEndTime)
		if err != nil {
			zap.L().Sugar().Info("get face entities err:", err.Error())
			goto END
		}
		zap.L().Sugar().Info("获取人脸抓拍个数：", len(faceEntities))
		zap.L().Sugar().Info("开始删除图片")

		for _, faceEntity := range faceEntities {
			<-deleteThreadChan
			wg.Add(1)
			go func(faceEntity *model.FaceTable) {
				defer func() {
					deleteThreadChan <- 1
					wg.Done()
				}()
				deleteWeedfs(client, faceEntity)
			}(faceEntity)

		}

		wg.Wait()
		zap.L().Sugar().Info("删除图片完成")
		zap.L().Sugar().Info("图片数量", picCount)
	END:
		totalCount += len(faceEntities)
		fmt.Println("~=~=~=~=~=~=~=~=~=~=~=~=")
		tmpEndTime = tmpStartTime
		tmpStartTime = tmpStartTime - 3600000
		if tmpStartTime < 0 {
			tmpStartTime = startTime
		}
		k++

	}

	wait := make(chan int)
	<-wait
}

//func deleteWeedfs(httpClient *http.Client, faceEntity *model.FaceTable) {
//	deleteWideImageUrls := []string{}
//	deleteCutboardImageUrls := []string{}
//
//	deleteWideImageUrls = append(deleteWideImageUrls, utils.ConverArceeURLToWeedUrl(faceEntity.ImageUri))
//	deleteCutboardImageUrls = append(deleteCutboardImageUrls, utils.ConverArceeURLToWeedUrl(faceEntity.ImageUri))
//
//	//fmt.Println("cut--->",deleteCutboardImageUrls)
//	for _, deleteWideImageUrl := range deleteWideImageUrls {
//		req, err := http.NewRequest("DELETE", deleteWideImageUrl, nil)
//		if err != nil {
//			fmt.Println("delete", deleteWideImageUrl, "err:", err.Error())
//			continue
//		}
//
//		ctx, cancel := context.WithCancel(context.TODO())
//		go func() {
//			<-time.After(2 * time.Second)
//			cancel()
//		}()
//
//		req.WithContext(ctx)
//		_, err = httpClient.Do(req)
//		if err != nil {
//			fmt.Println("delete", deleteWideImageUrl, "err:", err.Error())
//			continue
//		}
//
//	}
//
//	for _, deleteCutboardImageUrl := range deleteCutboardImageUrls {
//
//		// fmt.Println("cutboard:", deleteCutboardImageUrl)
//
//		req, err := http.NewRequest("DELETE", deleteCutboardImageUrl, nil)
//		if err != nil {
//			fmt.Println("delete", deleteCutboardImageUrl, "err:", err.Error())
//			continue
//		}
//
//		ctx, cancel := context.WithCancel(context.TODO())
//		go func() {
//			<-time.After(2 * time.Second)
//			cancel()
//		}()
//
//		req.WithContext(ctx)
//		_, err = httpClient.Do(req)
//		if err != nil {
//			fmt.Println("delete", deleteCutboardImageUrl, "err:", err.Error())
//			continue
//		}
//	}
//
//	// os.Exit(0)
//}

func deleteWeedfs(httpClient *http.Client, faceEntity *model.FaceTable) {
	if faceEntity.ImageUri != "" || faceEntity.CutboardImageUri != "" {
		deleteUrl(httpClient, faceEntity.ImageUri)
		deleteUrl(httpClient, faceEntity.CutboardImageUri)
	}
	//if faceEntity.CutboardImageUri != "" {
	//	req, err := http.NewRequest("DELETE", utils.ConverArceeURLToWeedUrl(faceEntity.CutboardImageUri), nil)
	//	if err != nil {
	//		fmt.Println("delete", faceEntity.CutboardImageUri, "err:", err.Error())
	//	}
	//
	//	ctx, cancel := context.WithCancel(context.TODO())
	//	go func() {
	//		<-time.After(2 * time.Second)
	//		cancel()
	//	}()
	//
	//	req.WithContext(ctx)
	//	resp, err := httpClient.Do(req)
	//	if resp != nil {
	//		defer resp.Body.Close()
	//	}
	//	if err != nil {
	//		fmt.Println("delete", faceEntity.ImageUri, "err:", err.Error())
	//		resp, _ = httpClient.Do(req)
	//	}
	//	if resp.StatusCode == http.StatusAccepted {
	//		zap.L().Debug("删除成功")
	//	}
	//	if resp.StatusCode == http.StatusNotFound {
	//		zap.L().Debug("图片不存在或者已经删除")
	//	}
	//	atomic.AddInt64(&picCount, 1)
	//}

}

func deleteUrl(httpClient *http.Client, url string) {
	if url == "" {
		return
	}
	weedUrl := utils.ConverArceeURLToWeedUrl(url)
	if weedUrl == "" {
		return
	}
	req, err := http.NewRequest("DELETE", weedUrl, nil)
	if err != nil {
		fmt.Println("delete", url, "err:", err.Error())
	}
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		<-time.After(global.HttpTimeOut * time.Second)
		cancel()
	}()

	req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Println("delete", url, "err:", err.Error())
		return
	}
	if resp == nil {
		fmt.Println("******************")
	}
	if resp.StatusCode == http.StatusAccepted {
		zap.L().Debug("删除成功")
	}
	if resp.StatusCode == http.StatusNotFound {
		zap.L().Debug("图片不存在或者已经删除")
	}
	atomic.AddInt64(&picCount, 1)
}

func weedShrink(httpClient *http.Client) error {
	s := `
   等执行完毕后，请手动在命令行,开启tmux执行下面的curl命令
	如果删除的文件有很多，会卡主很长时间，不要手动关闭
   curl http://%s:9333/vol/vacuum?garbageThreshold=%v

`
	//fmt.Println(1234,s)
	fmt.Printf(s, global.PicIP, global.ThresHold)
	//req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:9333/vol/vacuum?garbageThreshold=%v",
	//	global.PicIP, global.ThresHold), nil)
	//if err != nil {
	//	return err
	//}
	//
	//ctx, cancel := context.WithCancel(context.TODO())
	//go func() {
	//	<-time.After(100 * time.Second)
	//	cancel()
	//}()
	//
	//req.WithContext(ctx)
	//resp, err := httpClient.Do(req)
	//if err != nil {
	//	return err
	//}
	//zap.L().Sugar().Info(resp)

	return nil
}
