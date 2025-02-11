import React from "react";

import { Box, Table, TableCell, TableRow } from "@mui/material";

export const ProfileDetails: React.FC = () => {
  return (
    <Box>
      <Table>
        <TableRow>
          <TableCell>relacionamento:</TableCell>
          <TableCell>solteiro</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>interesses:</TableCell>
          <TableCell>amigos</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>sobre mim:</TableCell>
          <TableCell>Lorem ipsum dolor sit amet</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>etnia:</TableCell>
          <TableCell>asiático</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>religião:</TableCell>
          <TableCell>Budista</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>humor:</TableCell>
          <TableCell>amigável</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>moda:</TableCell>
          <TableCell>contemporânea</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>mora com:</TableCell>
          <TableCell>amigos</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>cidade natal:</TableCell>
          <TableCell>Piracicaba</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>página:</TableCell>
          <TableCell>www.flogão.com.br</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>paixões:</TableCell>
          <TableCell>tecnologia</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>esportes:</TableCell>
          <TableCell>frisbee</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>atividades:</TableCell>
          <TableCell>pular corda</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>livros:</TableCell>
          <TableCell>Duna</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>música:</TableCell>
          <TableCell>heteropoc</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>shows de tv:</TableCell>
          <TableCell>big brother</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>filmes:</TableCell>
          <TableCell>Matrix</TableCell>
        </TableRow>
        <TableRow>
          <TableCell>culinária:</TableCell>
          <TableCell>mexicana</TableCell>
        </TableRow>
      </Table>
    </Box>
  );
};
