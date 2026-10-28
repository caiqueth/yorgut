import { CommunityCategory } from "api/generated/model";

// Portuguese names of the community categories (Orkut's list). The API sends
// keys; a new category in the API fails the typecheck here until it gets a name.
export const communityCategoryLabels: Record<CommunityCategory, string> = {
  activities: "Atividades",
  alumni_schools: "Ex-alunos e Escolas",
  arts_entertainment: "Artes e Entretenimento",
  automotive: "Automotivo",
  business: "Negócios",
  cities_neighborhoods: "Cidades e Bairros",
  companies: "Empresas",
  computers_internet: "Computadores e Internet",
  countries_regional: "Países e Regiões",
  cultures_community: "Culturas e Comunidade",
  family_home: "Família e Lar",
  fashion_beauty: "Moda e Beleza",
  food_drink_wine: "Comidas, Bebidas e Vinhos",
  gay_lesbian_bi: "Gays, Lésbicas e Bi",
  games: "Jogos",
  government_politics: "Governo e Política",
  health_wellness_fitness: "Saúde, Bem-estar e Fitness",
  hobbies_crafts: "Hobbies e Trabalhos Manuais",
  individuals: "Pessoas",
  music: "Música",
  other: "Outros",
  pets_animals: "Animais de Estimação",
  recreation_sports: "Esportes e Lazer",
  religion_beliefs: "Religiões e Crenças",
  romance_relationships: "Romance e Relacionamentos",
  schools_education: "Escolas e Educação",
  science_history: "Ciência e História",
  travel: "Viagens",
};

export const communityCategoryLabel = (category: CommunityCategory) =>
  communityCategoryLabels[category] ?? category;

// Every category, by name, with "Outros" last.
export const communityCategoryOptions: CommunityCategory[] = [
  ...(Object.keys(communityCategoryLabels) as CommunityCategory[])
    .filter((category) => category !== CommunityCategory.CategoryOther)
    .sort((a, b) => communityCategoryLabels[a].localeCompare(communityCategoryLabels[b], "pt-BR")),
  CommunityCategory.CategoryOther,
];
