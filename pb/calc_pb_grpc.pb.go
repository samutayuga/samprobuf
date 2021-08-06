// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PrimerCalculatorClient is the client API for PrimerCalculator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrimerCalculatorClient interface {
	Calculate(ctx context.Context, in *CalculationRequest, opts ...grpc.CallOption) (*CalculationResponse, error)
}

type primerCalculatorClient struct {
	cc grpc.ClientConnInterface
}

func NewPrimerCalculatorClient(cc grpc.ClientConnInterface) PrimerCalculatorClient {
	return &primerCalculatorClient{cc}
}

func (c *primerCalculatorClient) Calculate(ctx context.Context, in *CalculationRequest, opts ...grpc.CallOption) (*CalculationResponse, error) {
	out := new(CalculationResponse)
	err := c.cc.Invoke(ctx, "/pb.PrimerCalculator/Calculate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PrimerCalculatorServer is the server API for PrimerCalculator service.
// All implementations must embed UnimplementedPrimerCalculatorServer
// for forward compatibility
type PrimerCalculatorServer interface {
	Calculate(context.Context, *CalculationRequest) (*CalculationResponse, error)
	mustEmbedUnimplementedPrimerCalculatorServer()
}

// UnimplementedPrimerCalculatorServer must be embedded to have forward compatible implementations.
type UnimplementedPrimerCalculatorServer struct {
}

func (UnimplementedPrimerCalculatorServer) Calculate(context.Context, *CalculationRequest) (*CalculationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Calculate not implemented")
}
func (UnimplementedPrimerCalculatorServer) mustEmbedUnimplementedPrimerCalculatorServer() {}

// UnsafePrimerCalculatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PrimerCalculatorServer will
// result in compilation errors.
type UnsafePrimerCalculatorServer interface {
	mustEmbedUnimplementedPrimerCalculatorServer()
}

func RegisterPrimerCalculatorServer(s grpc.ServiceRegistrar, srv PrimerCalculatorServer) {
	s.RegisterService(&PrimerCalculator_ServiceDesc, srv)
}

func _PrimerCalculator_Calculate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalculationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrimerCalculatorServer).Calculate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.PrimerCalculator/Calculate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrimerCalculatorServer).Calculate(ctx, req.(*CalculationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PrimerCalculator_ServiceDesc is the grpc.ServiceDesc for PrimerCalculator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PrimerCalculator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.PrimerCalculator",
	HandlerType: (*PrimerCalculatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Calculate",
			Handler:    _PrimerCalculator_Calculate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "calc_pb.proto",
}
