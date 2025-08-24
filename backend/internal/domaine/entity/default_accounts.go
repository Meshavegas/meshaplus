package entity

// DefaultAccounts représente les comptes par défaut créés pour chaque nouvel utilisateur
var DefaultAccounts = []Account{
	{Name: "Portefeuille (Cash)", Type: "cash", Currency: "XAF", Icon: "ion:cash", Color: "#00B894", Balance: 0, AccountNumber: nil},
	{Name: "Compte Bancaire", Type: "checking", Currency: "XAF", Icon: "fa5:university", Color: "#0984E3", Balance: 0, AccountNumber: nil},
	{Name: "MOMO", Type: "mobile_money", Currency: "XAF", Icon: "mci:cellphone", Color: "#FDCB6E", Balance: 0, AccountNumber: nil},
	{Name: "OM", Type: "mobile_money", Currency: "XAF", Icon: "mci:cellphone", Color: "#E17055", Balance: 0, AccountNumber: nil},
	{Name: "Épargne", Type: "savings", Currency: "XAF", Icon: "fa5:piggy-bank", Color: "#A29BFE", Balance: 0, AccountNumber: nil},
}
