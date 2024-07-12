package main

import (
	"context"       // Package pour la gestion des contextes
	"encoding/json" // Package pour l'encodage/décodage JSON
	"io/ioutil"     // Package pour les opérations d'entrée/sortie de fichier
	"log"           // Package pour la gestion des journaux
	"os"            // Package pour les opérations du système d'exploitation

	// Importation des définitions gRPC et des structures opérationnelles locales
	operation "rt0805/tp_app/operation"
	pb "rt0805/tp_app/operation_grpc"

	// Package pour l'utilisation de gRPC
	"google.golang.org/grpc"

	// Packages pour la gestion des chemins de fichiers, des entrées/sorties en buffer, des formats de chaîne
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	address  = "localhost:50051" // Adresse du serveur gRPC
	basePath = "../donnees/"     // Chemin de base pour les fichiers JSON
)

// convertToDevicePB convertit un objet operation.Device en pb.Device pour l'envoi via gRPC.
func convertToDevicePB(device *operation.Device) *pb.Device {
	// Crée une slice pour stocker les opérations converties
	operations := make([]*pb.Operation, len(device.Operations))
	
	// Parcourt chaque opération de l'appareil et la convertit en pb.Operation
	for i, op := range device.Operations {
		operations[i] = &pb.Operation{
			Type:         op.Type,         // Type d'opération
			HasSucceeded: op.HasSucceeded, // Statut de réussite de l'opération
		}
	}
	
	// Retourne un objet pb.Device avec le nom de l'appareil et les opérations converties
	return &pb.Device{
		Name:       device.Name,   // Nom de l'appareil
		Operations: operations,    // Liste des opérations converties
	}
}

func main() {
	// Créer un lecteur pour lire à partir de la console standard
	reader := bufio.NewReader(os.Stdin)

	// Demander le nom du fichier JSON à l'utilisateur
	fmt.Print("Entrez le nom du fichier JSON: ")
	
	// Lire l'entrée de l'utilisateur (nom du fichier JSON)
	fileName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Erreur lors de la lecture de l'entrée : %v", err)
	}

	// Nettoyer l'entrée du chemin pour enlever les espaces superflus et les nouvelles lignes
	fileName = strings.TrimSpace(fileName)

	// Construire le chemin complet du fichier en utilisant le chemin de base et le nom du fichier
	fullPath := filepath.Join(basePath, fileName)

	// Vérifier si le fichier existe
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Fatalf("Le fichier %s n'existe pas.", fullPath)
	}

	// Ouvrir le fichier JSON
	file, err := os.Open(fullPath)
	if err != nil {
		log.Fatalf("Impossible d'ouvrir le fichier JSON: %v", err)
	}
	// Assurer la fermeture du fichier après utilisation
	defer file.Close()

	// Lire le contenu du fichier JSON
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier JSON: %v", err)
	}

	// Initialiser une connexion gRPC avec le serveur
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Impossible de se connecter: %v", err)
	}
	// Assurer la fermeture de la connexion après utilisation
	defer conn.Close()

	// Créer un client gRPC pour interagir avec le serveur
	c := pb.NewDeviceServiceClient(conn)

	// Convertir les données JSON en objets Go (slice d'objets Device)
	var devices []operation.Device
	err = json.Unmarshal(data, &devices)
	if err != nil {
		log.Fatalf("Erreur lors de la conversion des données JSON: %v", err)
	}

	// Envoyer chaque appareil au serveur via gRPC
	for _, device := range devices {
		// Convertir l'appareil en un objet pb.Device compatible avec gRPC
		pbDevice := convertToDevicePB(&device)
		
		// Envoyer l'appareil converti au serveur via une requête gRPC
		_, err = c.SendData(context.Background(), &pb.DeviceDataRequest{Device: pbDevice})
		if err != nil {
			log.Fatalf("Erreur lors de l'envoi des données: %v", err)
		}
		// Journaliser le succès de l'envoi des données
		log.Printf("Données envoyées avec succès au serveur pour l'appareil %s.", device.Name)
	}
}
