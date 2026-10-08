import { Box, Pagination, Stack, TextField, Typography } from "@mui/material";
import { useSearchCommunities } from "api/community/queries";
import { useJoinCommunity } from "api/me/queries";
import { Button } from "components/atoms/Button";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React, { useState } from "react";

const PAGE_SIZE = 10;

interface IProps { }
export const DiscoverCommunities: React.FC<IProps> = () => {
  const [searchInput, setSearchInput] = useState("");
  // The term actually searched; only changes on submit.
  const [term, setTerm] = useState("");
  const [page, setPage] = useState(1);

  const { data, isLoading, isFetching } = useSearchCommunities(term, page, PAGE_SIZE);
  const { mutate: joinCommunity, isPending: isJoining, variables: joining } = useJoinCommunity();

  const handleSearch = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setTerm(searchInput.trim());
    setPage(1);
  };

  const count = data?.count ?? 0;
  const pageCount = Math.ceil(count / PAGE_SIZE);

  // Joining the last community of the last page would leave it empty.
  React.useEffect(() => {
    if (pageCount > 0 && page > pageCount) setPage(pageCount);
  }, [page, pageCount]);

  const items = (data?.communities ?? []).map((c) => (
    <Stack
      key={c.id}
      direction={"row"}
      justifyContent={"space-between"}
      alignItems={"center"}
      gap={1}
    >
      <PicAndTextHorizontal
        title={c.name}
        text={`${c.category} · ${c.memberCount} ${c.memberCount === 1 ? "membro" : "membros"}`}
        url={`/community/${c.id}`}
        imageSrc={c.profilePicture ?? ""}
        imageSize={64}
      />
      <Button
        onClick={() => joinCommunity({ communityId: c.id })}
        loading={isJoining && joining?.communityId === c.id}
      >
        participar
      </Button>
    </Stack>
  ));

  return (
    <Stack direction={"column"} gap={1}>
      <form onSubmit={handleSearch}>
        <Typography>Descubra novas comunidades</Typography>
        <Stack direction={"row"} alignItems={"center"} gap={1}>
          <TextField
            size={"small"}
            placeholder={"nome da comunidade"}
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
          />
          <Button type={"submit"}>Pesquisar</Button>
        </Stack>
      </form>
      <BoxWithTitle
        title={term ? `resultados para "${term}" (${count})` : `comunidades que você ainda não participa (${count})`}
        isLoading={isFetching && !isLoading}
      >
        {isLoading || items.length > 0
          ? <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />
          : (
            <Box sx={{ p: 1 }}>
              <Typography>
                {term ? "Nenhuma comunidade encontrada." : "Você já participa de todas as comunidades!"}
              </Typography>
            </Box>
          )}
        {pageCount > 1 && (
          <Stack direction={"row"} justifyContent={"center"} sx={{ p: 1 }}>
            <Pagination
              size={"small"}
              count={pageCount}
              page={page}
              onChange={(_e, value) => setPage(value)}
            />
          </Stack>
        )}
      </BoxWithTitle>
    </Stack>
  );
};
