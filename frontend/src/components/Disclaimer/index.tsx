import { Stack, Typography } from "@mui/material";
import React from "react";
import { Link } from "react-router-dom";

interface IProps {
  // "short": one line for the footer of every page, linking to /about.
  // "full": the whole text (login, about and help pages).
  variant: "short" | "full";
  // Smaller, lighter text for the full variant, when it sits next to other
  // text (the login page).
  muted?: boolean;
}
// Says that Yorgut is a personal, non-commercial project with no relation to
// Google or Orkut.
export const Disclaimer: React.FC<IProps> = ({ variant, muted }) => {
  if (variant === "short") {
    return (
      <Typography variant={"caption"} color={"text.secondary"} textAlign={"center"}>
        Yorgut é um projeto pessoal, sem fins lucrativos, sem relação com o Google ou o Orkut.{" "}
        <Link to={"/about"}>Saiba mais</Link>.
      </Typography>
    );
  }

  const textProps = muted
    ? { variant: "caption" as const, color: "text.secondary" }
    : { variant: "body2" as const };

  return (
    <Stack direction={"column"} gap={1}>
      <Typography {...textProps}>
        O Yorgut é um projeto pessoal, feito por hobby e sem fins lucrativos, inspirado no antigo
        Orkut. Não tem nenhuma relação com o Google LLC nem com o Orkut, e não é endossado,
        patrocinado ou mantido por eles.
      </Typography>
      <Typography {...textProps}>
        "Google" e "Orkut" são marcas do Google LLC. O login com Google serve apenas para identificar
        você; nenhum dado é compartilhado com terceiros.
      </Typography>
    </Stack>
  );
};
