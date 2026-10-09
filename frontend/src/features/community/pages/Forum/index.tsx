import { useCommunityTopics } from "api/community/queries";
import { Button } from "components/Button";
import { PicAndTextHorizontal } from "components/PicAndTextHorizontal";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import React from "react";
import { useNavigate, useParams } from "react-router-dom";
import { MembersOnly } from "features/community/components/MembersOnly";

interface IProps { }
export const Forum: React.FC<IProps> = () => {
  const params = useParams();
  const communityId = parseInt(params.communityId ?? "");

  const { data: topics, isLoading } = useCommunityTopics(communityId);
  const navigate = useNavigate();

  const handleGoToNewTopic = () => {
    navigate(`/community/${communityId}/forum/new-topic`)
  }

  const items = topics?.topics.map((t) => (
    <PicAndTextHorizontal
      title={t.title}
      text={t.lastReplyPreview}
      url={`topic/${t.id}`}
      imageSrc={t.authorProfilePicture}
    />
  )) ?? [];

  return (
    <BoxWithTitle title={"fórum"}>
      <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />
      <MembersOnly communityId={communityId} action={"criar tópicos"}>
        <Button onClick={handleGoToNewTopic}>+ novo tópico</Button>
      </MembersOnly>
    </BoxWithTitle>
  )
}
