import React from "react";
import { palette } from "theme";

import { Box, Paper, styled, Table, TableCell, TableRow } from "@mui/material";

const StyledTableRow = styled(TableRow)(({ theme }) => ({
  p: 0,
  m: 0,
  "&:nth-of-type(odd)": {
    backgroundColor: palette.lightBlue,
  },
  "&:nth-of-type(even)": {
    backgroundColor: "white",
  },
}));

const StyledTableCell = styled(TableCell)(({ theme }) => ({
  border: "none",
  padding: "0px 6px 0px 6px",
  "&:nth-of-type(odd)": {
    textAlign: "right",
    verticalAlign: "top",
  },
  "&:nth-of-type(even)": {
    textAlign: "left",
  },
}));

export const ProfileDetails: React.FC = () => {
  return (
    <Paper>
      <Box sx={{ width: "100%" }}>
        <Table
          size="small"
          sx={{
            width: "100%",
          }}
        >
          <StyledTableRow>
            <StyledTableCell>quem sou eu:</StyledTableCell>
            <StyledTableCell>
              Olá! Meu nome é Frizz, um majestoso gato tigrado de pelos macios e bigodes imponentes. Sou um verdadeiro
              mestre das artes felinas, especialista em cochilos estratégicos, corridas furtivas pela casa e, é claro,
              brincadeiras com cordinhas - meu passatempo favorito! Se você balançar uma na minha frente, prepare-se,
              porque eu vou caçar como um destemido leão da savana (ou quase isso). 🐾✨ Apesar da minha energia felina
              incontrolável, sou um gato extremamente bonzinho. Gosto de me esfregar nas pernas dos humanos para
              demonstrar afeto (e, talvez, para marcar território também, porque, né, segurança em primeiro lugar!).
              Adoro receber carinhos na cabeça e ronrono tão alto que poderia ser confundido com um motorzinho em
              funcionamento. Adoro uma coçadinha na barriguinha! Minha rotina é simples, mas nobre: dormir nos lugares
              mais inusitados, patrulhar os cômodos para garantir que nada tenha se movido sozinho e correr feito um
              foguete do nada - só para manter os humanos alertas. Se um dia você cruzar comigo, já sabe: traga uma
              cordinha, ofereça um carinho (com cautela) e aceite minha superioridade felina.
            </StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>relacionamento:</StyledTableCell>
            <StyledTableCell>solteiro</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>interesses:</StyledTableCell>
            <StyledTableCell>amigos</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>sobre mim:</StyledTableCell>
            <StyledTableCell>Lorem ipsum dolor sit amet</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>etnia:</StyledTableCell>
            <StyledTableCell>asiático</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>religião:</StyledTableCell>
            <StyledTableCell>Budista</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>humor:</StyledTableCell>
            <StyledTableCell>amigável</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>moda:</StyledTableCell>
            <StyledTableCell>contemporânea</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>mora com:</StyledTableCell>
            <StyledTableCell>amigos</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>cidade natal:</StyledTableCell>
            <StyledTableCell>Piracicaba</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>página:</StyledTableCell>
            <StyledTableCell>www.flogão.com.br</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>paixões:</StyledTableCell>
            <StyledTableCell>tecnologia</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>esportes:</StyledTableCell>
            <StyledTableCell>frisbee</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>atividades:</StyledTableCell>
            <StyledTableCell>pular corda</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>livros:</StyledTableCell>
            <StyledTableCell>Duna</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>música:</StyledTableCell>
            <StyledTableCell>heteropoc</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>shows de tv:</StyledTableCell>
            <StyledTableCell>big brother</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>filmes:</StyledTableCell>
            <StyledTableCell>Matrix</StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell>culinária:</StyledTableCell>
            <StyledTableCell>mexicana</StyledTableCell>
          </StyledTableRow>
        </Table>
      </Box>
    </Paper>
  );
};
