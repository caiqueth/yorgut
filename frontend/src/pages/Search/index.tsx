import { Box, LinearProgress, Stack, TextField, Typography } from "@mui/material";
import { useUserSearchUsers } from "api/user/queries";
import { Button } from "components/atoms/Button";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React, { useState } from "react";

interface IProps { }
export const Search: React.FC<IProps> = () => {
  const [searchName, setSearchName] = useState("");
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);

  const handleChangeSearchName = (e: any) => {
    setSearchName(e.target.value);
  }

  const { data, isFetching, refetch } = useUserSearchUsers(searchName, page, pageSize);

  const items = data?.users?.map((u) => (
    <PicAndTextHorizontal
      title={u.firstName}
      text={`${u.friendCount} ${u.friendCount > 1 ? "amigos" : "amigo"}`}
      url={`/profile/${u.id}`}
      imageSrc={u.profilePicture}
    />
  )) ?? [];

  return (
    <Box sx={{ width: "100%" }}>
      <Stack spacing={1}>
        <Box>
          <Typography>Pesquisar por usuários</Typography>
          <Stack
            direction={"row"}
            alignItems={"center"}
            gap={1}
            sx={{ width: "100%" }}
          >
            <TextField size={"small"} value={searchName} onChange={handleChangeSearchName} />
            <Button onClick={() => refetch()}>Pesquisar</Button>
          </Stack>
        </Box>
        <BoxWithTitle title={"resultados"}>
          <LinearProgress sx={{ visibility: isFetching ? "visible" : "hidden" }} />
          <StripedGrid items={items} elementsPerRow={1} isLoading={isFetching} />
        </BoxWithTitle>
      </Stack>
    </Box>
  );
};
