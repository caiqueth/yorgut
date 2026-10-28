import { Box, Stack, TextField, Typography } from "@mui/material";
import { useSearchCommunities } from "api/community/queries";
import { useJoinCommunity } from "api/me/queries";
import { Button } from "components/Button";
import { PicAndTextHorizontal } from "components/PicAndTextHorizontal";
import { CommunityRowSkeleton } from "components/CommunityRow";
import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { communityCategoryLabel } from "helpers/communityCategories";
import { ListPagination } from "components/ListPagination";
import React, { useState } from "react";

const PAGE_SIZE = 10;

export const DiscoverCommunities: React.FC = () => {
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
        text={`${c.private ? "🔒 " : ""}${communityCategoryLabel(c.category)} · ${c.memberCount} ${c.memberCount === 1 ? "membro" : "membros"}`}
        url={`/community/${c.id}`}
        imageSrc={c.profilePicture ?? ""}
        imageSize={64}
      />
      <Button
        onClick={() => joinCommunity({ communityId: c.id })}
        loading={isJoining && joining?.communityId === c.id}
        disabled={c.pending}
      >
        {c.pending ? "pedido enviado" : c.private ? "pedir para participar" : "participar"}
      </Button>
    </Stack>
  ));

  return (
    // The search goes inside the box: above it, it would sit between the
    // tabs and their content.
    <BoxWithTitle
      title={
        <>
          {term ? `resultados para "${term}"` : "comunidades que você ainda não participa"}{" "}
          <Count value={data?.count} />
        </>
      }
      showProgress={isFetching && !isLoading}
    >
      <Box component={"form"} onSubmit={handleSearch} sx={{ p: 1 }}>
        <Stack direction={"row"} alignItems={"center"} gap={1}>
          <Typography>Descubra novas comunidades:</Typography>
          <TextField
            placeholder={"nome da comunidade"}
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
          />
          <Button type={"submit"}>Pesquisar</Button>
        </Stack>
      </Box>
      {isLoading || items.length > 0
        ? <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={<CommunityRowSkeleton />} />
        : (
          <Box sx={{ p: 1 }}>
            <Typography>
              {term ? "Nenhuma comunidade encontrada." : "Você já participa de todas as comunidades!"}
            </Typography>
          </Box>
        )}
      {/* Joining the last community of the last page moves back a page. */}
      <ListPagination page={page} pageSize={PAGE_SIZE} count={data?.count} onChange={setPage} />
    </BoxWithTitle>
  );
};
