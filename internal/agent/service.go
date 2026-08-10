package agent

import (
	"context"
	"log"

	pb "github.com/mabuhasna8/IntelliOps/proto/automation/v1"
)

type Service struct {
	pb.UnimplementedAgentServiceServer
}

func NewService() *Service {
	return &Service{}
}

// WorkStream is a bidirectional stream.
// For now we just log messages and send no assignments.
func (s *Service) WorkStream(stream pb.AgentService_WorkStreamServer) error {
	log.Println("Agent connected")

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("agent stream closed: %v", err)
			return err
		}

		switch m := msg.Msg.(type) {
		case *pb.AgentMessage_Hello:
			log.Printf("hello from agent_id=%s env=%s", m.Hello.AgentId, m.Hello.EnvId)

		case *pb.AgentMessage_Request:
			log.Printf("work request: slots=%d", m.Request.Slots)
			// TODO: assign real tasks. For now, do nothing.
			// Example: stream.Send(&pb.CoreMessage{Msg: &pb.CoreMessage_Assignment{...}})

		case *pb.AgentMessage_Status:
			log.Printf("status: task=%s phase=%s msg=%s", m.Status.TaskId, m.Status.Phase, m.Status.Message)

		case *pb.AgentMessage_Log:
			log.Printf("log: task=%s line=%s", m.Log.TaskId, m.Log.Line)

		default:
			log.Printf("unknown agent message: %T", m)
		}
	}
}

// Example hook for scheduler later.
func (s *Service) AssignTask(ctx context.Context, task *pb.TaskSpec) error {
	// In a real impl, you'd push this onto some queue for streaming to agents.
	return nil
}

