package main

import (
	"testing" // Package pour les tests unitaires

	// Importation des définitions locales pour les opérations
	"rt0805/tp_app/operation"
)

// TestConvertToDevicePB teste la fonction convertToDevicePB
func TestConvertToDevicePB(t *testing.T) {
	// Créer un exemple de device
	device := operation.Device{
		Name: "testDevice", // Nom de l'appareil
		Operations: []operation.Operation{
			{Type: "CREATE", HasSucceeded: true},  // Opération de type CREATE réussie
			{Type: "UPDATE", HasSucceeded: false}, // Opération de type UPDATE échouée
		},
	}

	// Convertir le device en pb.Device
	pbDevice := convertToDevicePB(&device)

	// Vérifier si la conversion est correcte
	if pbDevice.Name != device.Name {
		t.Errorf("Nom incorrect, attendu: %s, obtenu: %s", device.Name, pbDevice.Name)
	}

	if len(pbDevice.Operations) != len(device.Operations) {
		t.Errorf("Nombre d'opérations incorrect, attendu: %d, obtenu: %d", len(device.Operations), len(pbDevice.Operations))
	}

	// Parcourir chaque opération et vérifier les détails
	for i, op := range device.Operations {
		// Vérifier si le type d'opération est correct
		if pbDevice.Operations[i].Type != op.Type {
			t.Errorf("Type d'opération incorrect pour l'opération %d, attendu: %s, obtenu: %s", i, op.Type, pbDevice.Operations[i].Type)
		}
		// Vérifier si le statut de réussite de l'opération est correct
		if pbDevice.Operations[i].HasSucceeded != op.HasSucceeded {
			t.Errorf("Statut de l'opération incorrect pour l'opération %d, attendu: %t, obtenu: %t", i, op.HasSucceeded, pbDevice.Operations[i].HasSucceeded)
		}
	}
}
