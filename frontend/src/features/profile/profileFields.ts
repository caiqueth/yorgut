// The descriptive profile fields, in display order. Shared by the profile
// details table and the edit page so both list the same fields and labels.
export type ProfileFieldKey =
  | "description"
  | "relationshipStatus"
  | "interests"
  | "about"
  | "religion"
  | "humor"
  | "fashion"
  | "livesWith"
  | "birthCity"
  | "website"
  | "passions"
  | "sports"
  | "hobbies"
  | "books"
  | "music"
  | "movies"
  | "cuisine";

export const profileFields: { key: ProfileFieldKey; label: string; multiline?: boolean }[] = [
  { key: "description", label: "quem sou eu:", multiline: true },
  { key: "relationshipStatus", label: "relacionamento:" },
  { key: "interests", label: "interesses:" },
  { key: "about", label: "sobre mim:", multiline: true },
  { key: "religion", label: "religião:" },
  { key: "humor", label: "humor:" },
  { key: "fashion", label: "moda:" },
  { key: "livesWith", label: "mora com:" },
  { key: "birthCity", label: "cidade natal:" },
  { key: "website", label: "página:" },
  { key: "passions", label: "paixões:", multiline: true },
  { key: "sports", label: "esportes:" },
  { key: "hobbies", label: "atividades:" },
  { key: "books", label: "livros:" },
  { key: "music", label: "música:" },
  { key: "movies", label: "filmes:" },
  { key: "cuisine", label: "culinária:" },
];
