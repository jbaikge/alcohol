package main

import (
	"log"
)

const IgnoreType = "IGNORE"

var types = map[string]string{
	"APPLE BRANDY":                             "Brandy",
	"BLENDED WHISKEY":                          "Whiskey",
	"BOTTLED IN BOND CORN WHISKEY":             "Whiskey",
	"BOTTLED IN BOND RYE WHISKEY":              "Whiskey",
	"BOTTLED IN BOND WHISKEY":                  "Whiskey",
	"BRANDY IMPORTED":                          "Brandy",
	"CANADIAN WHISKY":                          "Whiskey",
	"CHERRY BRANDY":                            "Brandy",
	"COCKTAILS (DOMESTIC)":                     IgnoreType,
	"COCKTAILS (IMPORTED)":                     IgnoreType,
	"COGNAC/ARMAGNAC":                          IgnoreType,
	"CORDIALS-LIQUEURS-SPECIALTIES (DOMESTIC)": IgnoreType,
	"CORDIALS-LIQUEURS-SPECIALTIES (IMPORTED)": IgnoreType,
	"CORN WHISKEY":                             "Whiskey",
	"EGG NOG":                                  IgnoreType,
	"FLAVORED BRANDY":                          "Brandy",
	"FLAVORED GIN":                             "Gin",
	"FLAVORED VODKA":                           IgnoreType,
	"GIN (DOMESTIC)":                           "Gin",
	"GIN (IMPORTED)":                           "Gin",
	"GRAPE BRANDY (DOMESTIC)":                  "Brandy",
	"IRISH WHISKEY":                            "Whiskey",
	"MISCELLANEOUS":                            IgnoreType,
	"MIXERS (NON ALCOHOLIC)":                   IgnoreType,
	"MOONSHINE":                                "Moonshine",
	"OTHER (DOMESTIC) BRANDY":                  "Brandy",
	"OTHER (IMPORTED) WHISKEY":                 "Whiskey",
	"PEACH BRANDY":                             "Brandy",
	"Rimmers":                                  IgnoreType,
	"ROCK & RYE":                               IgnoreType,
	"RUM (DOMESTIC)":                           "Rum",
	"RUM (IMPORTED)":                           "Rum",
	"SAMPLES":                                  IgnoreType,
	"SCOTCH WHISKY":                            "Scotch",
	"SLOE GIN":                                 "Gin",
	"SPECIALTY BOTTLES   (DOMESTIC)":           "Specialty",
	"SPECIALTY BOTTLES (IMPORTED)":             "Specialty",
	"STRAIGHT BOURBON WHISKEY":                 "Bourbon",
	"STRAIGHT RYE WHISKEY":                     "Whiskey",
	"TENNESSEE WHISKEY":                        "Whiskey",
	"TEQUILA":                                  "Tequila",
	"VERMOUTH (DOMESTIC)":                      IgnoreType,
	"VERMOUTH (IMPORTED)":                      IgnoreType,
	"VIRGINIA FRUIT WINES":                     IgnoreType,
	"VIRGINIA MISCELLANEOUS WINE":              IgnoreType,
	"VIRGINIA PINK WINE":                       IgnoreType,
	"VIRGINIA RED WINE":                        IgnoreType,
	"VIRGINIA RED TABLE WINE":                  IgnoreType,
	"VIRGINIA PINK TABLE WINE":                 IgnoreType,
	"VIRGINIA WHITE TABLE WINE":                IgnoreType,
	"VIRGINIA SPARKLING WINE":                  IgnoreType,
	"VIRGINIA WHITE WINE":                      IgnoreType,
	"VODKA (DOMESTIC)":                         IgnoreType,
	"VODKA (IMPORTED)":                         IgnoreType,
	"WHISKEY":                                  "Whiskey",
}

func Type(in string) (t string) {
	t, ok := types[in]
	if !ok {
		t = in
		log.Printf("Unknown type: %s", in)
	}
	return
}
