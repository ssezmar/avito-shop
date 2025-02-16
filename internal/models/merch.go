package models

import "fmt"

type Merch struct {
    Name  string `json:"name"`
    Price int    `json:"price"`
}

var merchList = []Merch{
    {"t-shirt", 80},
    {"cup", 20},
    {"book", 50},
    {"pen", 10},
    {"powerbank", 200},
    {"hoody", 300},
    {"umbrella", 200},
    {"socks", 10},
    {"wallet", 50},
    {"pink-hoody", 500},
}

func GetMerchList() []Merch {
    return merchList
}

func GetMerchByName(name string) (Merch, error) {
    for _, merch := range merchList {
        if merch.Name == name {
            return merch, nil
        }
    }
    return Merch{}, fmt.Errorf("Merch not found")
}