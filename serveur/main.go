package main

import (
	"context"                         // Package pour la gestion des contextes
	"log"                             // Package pour la gestion des logs
	"net"                             // Package pour la gestion des connexions réseau
	"rt0805/tp_app/operation"         // Importation des structures opérationnelles
	pb "rt0805/tp_app/operation_grpc" // Importation des définitions gRPC
	"strconv"                         // Package pour les conversions de types (par exemple, int en string)
	"time"                            // Package pour la gestion du temps

	"go.mongodb.org/mongo-driver/bson"          // Package pour utiliser BSON (Binary JSON)
	"go.mongodb.org/mongo-driver/mongo"         // Package pour interagir avec MongoDB
	"go.mongodb.org/mongo-driver/mongo/options" // Package pour les options MongoDB
	"google.golang.org/grpc"                    // Package pour utiliser gRPC
)

const (
    port             = ":50051"     // Port sur lequel le serveur gRPC écoute
    mongoDBHost      = "localhost"  // Adresse du serveur MongoDB
    mongoDBPort      = 27017        // Port du serveur MongoDB
    mongoDBUsername  = "root"       // Nom d'utilisateur MongoDB
    mongoDBPassword  = "root"       // Mot de passe MongoDB
)

// Définition du serveur gRPC
type server struct {
    pb.UnimplementedDeviceServiceServer
}

// Implémentation de la méthode SendData du service DeviceService
func (s *server) SendData(ctx context.Context, in *pb.DeviceDataRequest) (*pb.DeviceDataResponse, error) {
    // Construction de l'URI de connexion MongoDB
    mongoURI := "mongodb://" + mongoDBUsername + ":" + mongoDBPassword + "@" + mongoDBHost + ":" + strconv.Itoa(mongoDBPort)
    client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal(err)
    }

    // Création d'un contexte avec un timeout de 10 secondes
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Connexion à MongoDB
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)

    // Sélection de la collection MongoDB
    collection := client.Database("operation").Collection("devices")

    // Recherche de l'appareil existant dans la collection
    var existingDevice operation.Device
    filter := bson.M{"device_name": in.Device.Name}
    err = collection.FindOne(ctx, filter).Decode(&existingDevice)

    // Calcul du nombre total et du nombre d'échecs des opérations
    totalOps := len(in.Device.Operations)
    failedOps := 0
    for _, op := range in.Device.Operations {
        if !op.HasSucceeded {
            failedOps++
        }
    }

    // Si l'appareil n'existe pas, insertion d'un nouvel appareil
    if err != nil {
        _, err = collection.InsertOne(ctx, bson.M{
            "device_name":       in.Device.Name,
            "total_operations":  totalOps,
            "failed_operations": failedOps,
            "operations":        in.Device.Operations,
        })
        if err != nil {
            log.Fatalf("Erreur lors de la création de l'appareil: %v", err)
        }
    } else {
        // Mise à jour de l'appareil existant
        _, err = collection.UpdateOne(
            ctx,
            filter,
            bson.M{
                "$inc": bson.M{
                    "total_operations":  totalOps,
                    "failed_operations": failedOps,
                },
                "$push": bson.M{"operations": bson.M{"$each": in.Device.Operations}},
            },
        )
        if err != nil {
            log.Fatalf("Erreur lors de la mise à jour de l'appareil: %v", err)
        }
    }

    // Retourne une réponse de succès
    return &pb.DeviceDataResponse{Success: true}, nil
}

// Fonction principale pour initialiser et démarrer le serveur gRPC
func main() {
    // Écoute sur le port spécifié
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Échec de l'écoute: %v", err)
    }

    // Création d'un nouveau serveur gRPC
    s := grpc.NewServer()

    // Enregistrement du service DeviceService avec le serveur gRPC
    pb.RegisterDeviceServiceServer(s, &server{})

    // Démarrage du serveur gRPC
    log.Printf("Serveur gRPC démarré sur le port %s", port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Échec du serveur: %v", err)
    }
}
