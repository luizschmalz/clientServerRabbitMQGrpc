package main

import (
	"context"
	"log"
	"net"
	pb "cliente-servidor/proto"
	"cliente-servidor/imageutils"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedImageServiceServer
}

func (s *server) ConvertToGray(ctx context.Context, req *pb.ImageRequest) (*pb.ImageReply, error) {
	img, err := imageutils.BytesToImage(req.ImageData)
	if err != nil {
		return nil, err
	}
	grayImg := imageutils.ToGray(img)
	return &pb.ImageReply{ImageData: imageutils.ImageToBytes(grayImg)}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterImageServiceServer(grpcServer, &server{})
	log.Println("gRPC server running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
