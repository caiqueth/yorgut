import { useCommunityTopics } from "api/community/queries";
import { Button } from "components/atoms/Button";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";
import { useNavigate, useParams } from "react-router-dom";

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

  items.push(<Button onClick={handleGoToNewTopic}>+ novo tópico</Button>);

  return (
    <BoxWithTitle title={"fórum"}>
      <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />
    </BoxWithTitle>
  )
}
