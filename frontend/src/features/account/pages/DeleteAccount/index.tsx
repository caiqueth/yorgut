import { Box, Checkbox, FormControlLabel, Stack, Typography } from "@mui/material";
import { useDeleteMe } from "api/me/queries";
import { BoxWithTitle } from "components/BoxWithTitle";
import { Button } from "components/Button";
import { useLogout } from "hooks/useLogout";
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

const deleted = [
  "seu perfil, sua foto e as fotos do seu álbum;",
  "os recados e depoimentos que você escreveu e recebeu;",
  "suas amizades, pedidos de amizade, estrelas de fã e avaliações, dadas e recebidas;",
  "os tópicos que você criou no fórum, com todas as respostas deles, e as suas respostas em outros tópicos;",
  "as enquetes que você criou e os seus votos;",
  "sua participação em comunidades, seus convites e pedidos para participar, e o registro das suas visitas.",
];

// Asks before deleting the account, saying what goes and what happens to the
// communities the user owns. Deleting logs out.
export const DeleteAccount: React.FC = () => {
  const navigate = useNavigate();
  const logout = useLogout();
  const { mutate: deleteMe, isPending, error } = useDeleteMe();
  const [understood, setUnderstood] = useState(false);

  return (
    <BoxWithTitle title={"excluir conta"}>
      <Stack direction={"column"} gap={2} sx={{ p: 2 }} alignItems={"flex-start"}>
        <Typography variant={"body2"}>
          Excluir a sua conta apaga para sempre, sem como recuperar:
        </Typography>
        <Box component={"ul"} sx={{ m: 0, pl: 3 }}>
          {deleted.map((item) => (
            <Typography key={item} component={"li"} variant={"body2"}>{item}</Typography>
          ))}
        </Box>
        <Typography variant={"body2"}>
          As comunidades de que você é dono(a) passam para o membro mais antigo de cada uma. As que não
          têm outros membros são apagadas.
        </Typography>
        <FormControlLabel
          control={<Checkbox size={"small"} checked={understood} onChange={(e) => setUnderstood(e.target.checked)} />}
          label={<Typography variant={"body2"}>Entendo que isto não pode ser desfeito.</Typography>}
        />
        {error && <Typography variant={"body2"} color={"error"}>{error.message}</Typography>}
        <Stack direction={"row"} gap={1} flexWrap={"wrap"}>
          <Button onClick={() => navigate("/my-account")} disabled={isPending}>cancelar</Button>
          <Button disabled={!understood} loading={isPending} onClick={() => deleteMe(undefined, { onSuccess: logout })}>
            excluir minha conta para sempre
          </Button>
        </Stack>
      </Stack>
    </BoxWithTitle>
  );
};
