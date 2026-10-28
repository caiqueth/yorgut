package transport

// CommunityCategory is one of the fixed community categories (Orkut's list).
// The API stores and sends these keys; the frontend shows a Portuguese name
// for each (helpers/communityCategories.ts).
type CommunityCategory string

const (
	CategoryActivities            CommunityCategory = "activities"
	CategoryAlumniSchools         CommunityCategory = "alumni_schools"
	CategoryArtsEntertainment     CommunityCategory = "arts_entertainment"
	CategoryAutomotive            CommunityCategory = "automotive"
	CategoryBusiness              CommunityCategory = "business"
	CategoryCitiesNeighborhoods   CommunityCategory = "cities_neighborhoods"
	CategoryCompanies             CommunityCategory = "companies"
	CategoryComputersInternet     CommunityCategory = "computers_internet"
	CategoryCountriesRegional     CommunityCategory = "countries_regional"
	CategoryCulturesCommunity     CommunityCategory = "cultures_community"
	CategoryFamilyHome            CommunityCategory = "family_home"
	CategoryFashionBeauty         CommunityCategory = "fashion_beauty"
	CategoryFoodDrinkWine         CommunityCategory = "food_drink_wine"
	CategoryGayLesbianBi          CommunityCategory = "gay_lesbian_bi"
	CategoryGames                 CommunityCategory = "games"
	CategoryGovernmentPolitics    CommunityCategory = "government_politics"
	CategoryHealthWellnessFitness CommunityCategory = "health_wellness_fitness"
	CategoryHobbiesCrafts         CommunityCategory = "hobbies_crafts"
	CategoryIndividuals           CommunityCategory = "individuals"
	CategoryMusic                 CommunityCategory = "music"
	CategoryOther                 CommunityCategory = "other"
	CategoryPetsAnimals           CommunityCategory = "pets_animals"
	CategoryRecreationSports      CommunityCategory = "recreation_sports"
	CategoryReligionBeliefs       CommunityCategory = "religion_beliefs"
	CategoryRomanceRelationships  CommunityCategory = "romance_relationships"
	CategorySchoolsEducation      CommunityCategory = "schools_education"
	CategoryScienceHistory        CommunityCategory = "science_history"
	CategoryTravel                CommunityCategory = "travel"
)

// CommunityCategories lists every category, in the order above.
var CommunityCategories = []CommunityCategory{
	CategoryActivities, CategoryAlumniSchools, CategoryArtsEntertainment, CategoryAutomotive,
	CategoryBusiness, CategoryCitiesNeighborhoods, CategoryCompanies, CategoryComputersInternet,
	CategoryCountriesRegional, CategoryCulturesCommunity, CategoryFamilyHome, CategoryFashionBeauty,
	CategoryFoodDrinkWine, CategoryGayLesbianBi, CategoryGames, CategoryGovernmentPolitics,
	CategoryHealthWellnessFitness, CategoryHobbiesCrafts, CategoryIndividuals, CategoryMusic,
	CategoryOther, CategoryPetsAnimals, CategoryRecreationSports, CategoryReligionBeliefs,
	CategoryRomanceRelationships, CategorySchoolsEducation, CategoryScienceHistory, CategoryTravel,
}

// IsValid reports whether c is one of CommunityCategories.
func (c CommunityCategory) IsValid() bool {
	for _, category := range CommunityCategories {
		if c == category {
			return true
		}
	}
	return false
}
