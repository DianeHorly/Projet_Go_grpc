package main

import (
	"context"                         // Package pour la gestion des contextes
	pb "rt0805/tp_app/operation_grpc" // Importation des définitions gRPC
	"testing"                         // Package pour les tests unitaires

	"go.mongodb.org/mongo-driver/mongo"          // Package pour interagir avec MongoDB
	"go.mongodb.org/mongo-driver/mongo/options"  // Package pour les options MongoDB
	"go.mongodb.org/mongo-driver/mongo/readpref" // Package pour les préférences de lecture MongoDB
	// "google.golang.org/grpc"                  // Importation (commentée) de gRPC pour un usage potentiel
)

// TestSendData teste la méthode SendData du serveur gRPC
func TestSendData(t *testing.T) {
    // Créer un client MongoDB connecté à un serveur MongoDB local pour les tests
    ctx := context.Background()
    mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        t.Fatalf("Erreur de connexion à la base de données MongoDB: %v", err)
    }
    // Déconnexion du client MongoDB une fois le test terminé
    defer mongoClient.Disconnect(ctx)

    // Vérifier que le serveur MongoDB est disponible
    if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
        t.Fatalf("Le serveur MongoDB n'est pas disponible: %v", err)
    }

    // Créer une instance du serveur gRPC
    deviceServiceServer := &server{}

    // Préparer une requête avec des données de test pour envoyer au serveur
    request := &pb.DeviceDataRequest{
        Device: &pb.Device{
            Name: "testDevice", // Nom de l'appareil de test
            Operations: []*pb.Operation{
                {Type: "CREATE", HasSucceeded: true},  // Opération de type CREATE réussie
                {Type: "UPDATE", HasSucceeded: false}, // Opération de type UPDATE échouée
            },
        },
    }

    // Envoyer la requête au serveur
    response, err := deviceServiceServer.SendData(ctx, request)
    if err != nil {
        t.Fatalf("Erreur lors de l'envoi des données: %v", err)
    }

    // Vérifier la réponse du serveur
    if !response.Success {
        t.Error("La réponse du serveur indique un échec")
    }
}
