package runner

import (
	"context"
	"okapi/helpers/jobs"
	lib_cmd "okapi/lib/cmd"
	"okapi/lib/task"
	proto "okapi/protos/runner"
)

// Server runner gRPC server
type Server struct {
	proto.UnimplementedRunnerServer
}

// Execute run the task through gRPC
func (srv *Server) Execute(req *proto.Request, stream proto.Runner_ExecuteServer) error {
	job, ctx, err := jobs.FromRPC(req, lib_cmd.Context)

	if err != nil {
		stream.Send(&proto.Response{
			Status: Failed,
			Info:   err.Error(),
		})

		return err
	}

	err = task.Exec(job, &ctx)

	if err != nil {
		stream.Send(&proto.Response{
			Status: Failed,
			Info:   err.Error(),
		})

		return err
	}

	stream.Send(&proto.Response{
		Status: Success,
	})

	return nil
}

// Enqueue execute task without waiting for response
func (srv *Server) Enqueue(context context.Context, req *proto.Request) (*proto.Response, error) {
	job, ctx, err := jobs.FromRPC(req, lib_cmd.Context)
	res := new(proto.Response)

	if err != nil {
		res.Status = Failed
		return res, err
	}

	res.Status = Success
	go task.Exec(job, &ctx)
	return res, nil
}
