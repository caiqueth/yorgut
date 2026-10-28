import { useIdParam } from "hooks/useIdParam";
import { useCommunityTopics } from "api/community/queries";
import { Button } from "components/Button";
import { PicAndTextHorizontal, PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";
import { BoxWithTitle } from "components/BoxWithTitle";
import { InlineSkeleton } from "components/InlineSkeleton";
import { StripedGrid } from "components/StripedGrid";
import { ListPagination } from "components/ListPagination";
import { FORUM_PAGE_SIZE } from "helpers/paging";
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { MembersOnly } from "features/community/components/MembersOnly";
import { useForumAccess } from "features/community/useForumAccess";

export const Forum: React.FC = () => {
  const communityId = useIdParam("communityId");

  const { canRead, isLoading: isLoadingAccess } = useForumAccess(communityId);
  const [page, setPage] = useState(1);
  const { data: topics, isLoading: isLoadingData, isFetching } = useCommunityTopics(communityId, page, FORUM_PAGE_SIZE, canRead);
  // The request waits for the access check; both count as loading.
  const isLoading = isLoadingAccess || isLoadingData;
  const navigate = useNavigate();

  const handleGoToNewTopic = () => {
    navigate(`/community/${communityId}/forum/new-topic`)
  }

  const items = topics?.topics.map((t) => (
    <PicAndTextHorizontal
      key={t.id}
      title={t.title}
      titleSuffix={` (${t.replies} ${t.replies === 1 ? "resposta" : "respostas"})`}
      text={t.lastReplyPreview}
      url={`topic/${t.id}`}
      imageUrl={`/profile/${t.authorId}`}
      imageSrc={t.authorProfilePicture}
    />
  )) ?? [];

  if (!isLoadingAccess && !canRead) {
    return (
      <BoxWithTitle title={"fórum"}>
        <MembersOnly communityId={communityId} action={"ver o fórum desta comunidade privada"} />
      </BoxWithTitle>
    );
  }

  return (
    <BoxWithTitle title={<>fórum ({topics ? topics.count : <InlineSkeleton width={"2ch"} />} tópicos)</>} showProgress={isFetching && !isLoading}>
      <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={<PicAndTextHorizontalSkeleton lineWidths={[200, 160]} />} />
      <ListPagination page={page} pageSize={FORUM_PAGE_SIZE} count={topics?.count} onChange={setPage} />
      <MembersOnly communityId={communityId} action={"criar tópicos"}>
        <Button onClick={handleGoToNewTopic}>+ novo tópico</Button>
      </MembersOnly>
    </BoxWithTitle>
  )
}
