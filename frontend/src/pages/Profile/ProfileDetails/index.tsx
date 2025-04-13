import React, { use } from "react";

import { Box, Paper, Table } from "@mui/material";
import { StyledTableCell, StyledTableRow } from "./styles";
import { useUserProfile } from "api/user/queries";

interface IProfileDetailsProps {
  userId: number;
}
export const ProfileDetails: React.FC<IProfileDetailsProps> = ({ userId }) => {
  const { data: userProfile } = useUserProfile(userId);

  const data = [
    {
      label: "quem sou eu:",
      value: userProfile?.description,
    },
    {
      label: "relacionamento:",
      value: userProfile?.relationshipStatus,
    },
    {
      label: "interesses:",
      value: userProfile?.interests,
    },
    {
      label: "sobre mim:",
      value: userProfile?.about,
    },
    {
      label: "religião:",
      value: userProfile?.religion,
    },
    {
      label: "humor:",
      value: userProfile?.humor,
    },
    {
      label: "moda:",
      value: userProfile?.fashion,
    },
    {
      label: "mora com:",
      value: userProfile?.livesWith,
    },
    {
      label: "cidade natal:",
      value: userProfile?.birthCity,
    },
    {
      label: "página:",
      value: userProfile?.website,
    },
    {
      label: "paixões:",
      value: userProfile?.passions,
    },
    {
      label: "esportes:",
      value: userProfile?.sports,
    },
    {
      label: "atividades:",
      value: userProfile?.hobbies,
    },
    {
      label: "livros:",
      value: userProfile?.books,
    },
    {
      label: "música:",
      value: userProfile?.music,
    },
    {
      label: "filmes:",
      value: userProfile?.movies,
    },
    {
      label: "culinária:",
      value: userProfile?.cuisine,
    },
  ];

  return (
    <Paper>
      <Box sx={{ width: "100%" }}>
        <Table
          size="small"
          sx={{
            width: "100%",
          }}
        >
          {data.map(({ label, value }) => (
            <StyledTableRow>
              <StyledTableCell sx={{ width: 140 }}>{label}</StyledTableCell>
              <StyledTableCell>{value}</StyledTableCell>
            </StyledTableRow>
          ))}
        </Table>
      </Box>
    </Paper>
  );
};
