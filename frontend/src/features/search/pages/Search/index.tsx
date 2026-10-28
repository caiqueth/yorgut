import { Box, Stack, TextField, Typography } from "@mui/material";
import { useUserSearchUsers } from "api/user/queries";
import { Button } from "components/Button";
import { UserRow, UserRowSkeleton } from "components/UserRow";
import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { ListPagination } from "components/ListPagination";
import { StripedGrid } from "components/StripedGrid";
import { LIST_PAGE_SIZE } from "helpers/paging";
import React, { useState } from "react";

export const Search: React.FC = () => {
  const [searchInput, setSearchInput] = useState("");
  // The name actually searched; only changes on submit.
  const [searchName, setSearchName] = useState("");
  const [page, setPage] = useState(1);

  const { data, isLoading, isFetching } = useUserSearchUsers(searchName, page, LIST_PAGE_SIZE);

  const handleSearch = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSearchName(searchInput.trim());
    setPage(1);
  };

  const items = data?.users?.map((u) => <UserRow key={u.id} user={u} />) ?? [];

  return (
    <Box sx={{ width: "100%" }}>
      <Stack spacing={1}>
        <form onSubmit={handleSearch}>
          <Typography>Pesquisar por usuários</Typography>
          <Stack
            direction={"row"}
            alignItems={"center"}
            gap={1}
            sx={{ width: "100%" }}
          >
            <TextField
              placeholder={"nome"}
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
            />
            <Button type={"submit"}>Pesquisar</Button>
          </Stack>
        </form>
        <BoxWithTitle title={<>resultados <Count value={data?.count} /></>} showProgress={isFetching && !isLoading}>
          <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={<UserRowSkeleton />} />
          <ListPagination page={page} pageSize={LIST_PAGE_SIZE} count={data?.count} onChange={setPage} />
        </BoxWithTitle>
      </Stack>
    </Box>
  );
};
