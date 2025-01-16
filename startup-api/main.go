package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Startup struct with all fields from your MongoDB document
type Startup struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date                  int64              `bson:"date" json:"date"`
	StartupID             string             `bson:"startupId" json:"startupId"`
	UID                   string             `bson:"uid" json:"uid"`
	Financing             interface{}        `bson:"financing" json:"financing"`
	Type                  string             `bson:"type" json:"type"`
	HasMetrics            bool               `bson:"hasMetrics" json:"hasMetrics"`
	OpenToOffers          bool               `bson:"openToOffers" json:"openToOffers"`
	AskingPrice           interface{}        `bson:"askingPrice" json:"askingPrice"`
	RevenueMultiple       interface{}        `bson:"revenueMultiple" json:"revenueMultiple"`
	ProfitMultiple        interface{}        `bson:"profitMultiple" json:"profitMultiple"`
	AdvisorID             interface{}        `bson:"advisorId" json:"advisorId"`
	ListingHeadline       string             `bson:"listingHeadline" json:"listingHeadline"`
	ListingType           string             `bson:"listingType" json:"listingType"`
	ManagedByMicro        interface{}        `bson:"managedByMicro" json:"managedByMicro"`
	TopPickedAt           interface{}        `bson:"topPickedAt" json:"topPickedAt"`
	BusinessVerified      bool               `bson:"businessVerified" json:"businessVerified"`
	TotalRevenueAnnual    interface{}        `bson:"totalRevenueAnnual" json:"totalRevenueAnnual"`
	TotalProfitAnnual     interface{}        `bson:"totalProfitAnnual" json:"totalProfitAnnual"`
	TotalGrowthAnnual     interface{}        `bson:"totalGrowthAnnual" json:"totalGrowthAnnual"`
	Location              string             `bson:"location" json:"location"`
	DateFounded           int64              `bson:"dateFounded" json:"dateFounded"`
	Team                  string             `bson:"team" json:"team"`
	About                 string             `bson:"about" json:"about"`
	Status                string             `bson:"status" json:"status"`
	Revenue               interface{}        `bson:"revenue" json:"revenue"`
	Customers             string             `bson:"customers" json:"customers"`
	Keywords              interface{}        `bson:"keywords" json:"keywords"`
	AnnualProfit          interface{}        `bson:"annualProfit" json:"annualProfit"`
	GrowthAnnual          interface{}        `bson:"growthAnnual" json:"growthAnnual"`
	TechStack             string             `bson:"techStack" json:"techStack"`
	TechStackKeywords     interface{}        `bson:"techStackKeywords" json:"techStackKeywords"`
	BusinessModel         string             `bson:"businessModel" json:"businessModel"`
	BusinessModelKeywords interface{}        `bson:"businessModelKeywords" json:"businessModelKeywords"`
	Competitors           interface{}        `bson:"competitors" json:"competitors"`
	WeeklyViews           interface{}        `bson:"weeklyViews" json:"weeklyViews"`
	Score                 interface{}        `bson:"score" json:"score"`
}

var client *mongo.Client
var startupCollection *mongo.Collection

func main() {
	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("MONGODB_URL")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")

	startupCollection = client.Database("acquire").Collection("acquire")

	// Set up routes
	http.HandleFunc("/startup", handleStartup)
	http.HandleFunc("/startups", getAllStartups)

	fmt.Println("Server starting on port 5050...")
	if err := http.ListenAndServe(":5050", nil); err != nil {
		panic(err)
	}
}

// CRUD handlers
func handleStartup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		createStartup(w, r)
	case http.MethodGet:
		getStartup(w, r)
	case http.MethodPut:
		updateStartup(w, r)
	case http.MethodDelete:
		deleteStartup(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createStartup(w http.ResponseWriter, r *http.Request) {
	var startup Startup
	if err := json.NewDecoder(r.Body).Decode(&startup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startup.Date = time.Now().UnixMilli()

	result, err := startupCollection.InsertOne(context.Background(), startup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Startup created successfully",
		"id":      result.InsertedID,
	})
}

func getStartup(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var startup Startup
	err = startupCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&startup)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Startup not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(startup)
}

func getAllStartups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var startups []Startup
	cursor, err := startupCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &startups); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(startups)
}

func updateStartup(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var startup Startup
	if err := json.NewDecoder(r.Body).Decode(&startup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := bson.M{"$set": startup}
	result, err := startupCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		update,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Startup not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Startup updated successfully",
		"modified_count": result.ModifiedCount,
	})
}

func deleteStartup(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	result, err := startupCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Startup not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Startup deleted successfully",
	})
}
