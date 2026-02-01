package client

import (
	"fmt"

	"github.com/MediStatTech/auth-client/client_options"
	services_v1 "github.com/MediStatTech/auth-client/pb/go/services/v1"
	xdscreds "google.golang.org/grpc/credentials/xds"
	log "github.com/MediStatTech/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/xds" // Import XDS support for production service discovery
)

type Facade struct {
	conn *grpc.ClientConn
	log  *log.Logger

	// Position service client
	Position services_v1.PositionServiceClient
	Staffs   services_v1.StaffsServiceClient
	Jwt      services_v1.JwtServiceClient
}

func New(
	o *client_options.Options,
) (*Facade, error) {
	var target string

	// Check if custom address is provided
	if o.AddressName != "" {
		target = o.AddressName
	} else {
		// XDS service discovery for production
		// This will resolve to the service endpoint via xDS control plane
		target = "todo-service.svc.cluster.local:8443"
	}

	// Connection options
	dialOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // Load balancing
	}

	// For development, use insecure connection
	// For production with XDS, use XDS credentials with insecure fallback
	if o.ENV != nil && o.ENV.IsDev() {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else if o.AddressName != "" && (len(o.AddressName) < 6 || o.AddressName[:6] != "xds://") {
		// Custom non-XDS address, use insecure
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// For XDS addresses in non-dev environments, use XDS credentials
		creds, err := xdscreds.NewClientCredentials(xdscreds.ClientOptions{
			FallbackCreds: insecure.NewCredentials(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create xds credentials: %w", err)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}

	conn, err := grpc.NewClient(target, dialOpts...)
	if err != nil {
		o.Log.Error("Failed to connect to auth service", map[string]interface{}{
			"target": target,
			"error":  err.Error(),
		})
		return nil, err
	}

	o.Log.Info("Connected to auth service", map[string]interface{}{
		"target": target,
	})

	return &Facade{
		conn:     conn,
		log:      o.Log,
		Position: services_v1.NewPositionServiceClient(conn),
		Staffs:   services_v1.NewStaffsServiceClient(conn),
		Jwt:      services_v1.NewJwtServiceClient(conn),
	}, nil
}

// Close closes the gRPC connection
func (c *Facade) Close() error {
	return c.conn.Close()
}
