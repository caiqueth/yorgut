import React from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";

import { Box } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useMeAmFriend } from "api/me/queries";
import { useAddFriend, useRemoveFriend, useReplyFriendship } from "api/friendship/queries";

interface IProfileMenuProps {
  userId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ userId }) => {
  const navigate = useNavigate();
  const { data: friendshipState } = useMeAmFriend(userId);
  const { mutate: replyFriendship } = useReplyFriendship();
  const { mutate: removeFriend } = useRemoveFriend(userId);
  const { mutate: addFriend } = useAddFriend();

  const handleDoNothing = () => { };

  const handleAddFriend = () => {
    addFriend({
      userRequesteeId: userId,
    })
  };

  const handleRemoveFriend = () => {
    if (!friendshipState?.friendshipId) return;
    removeFriend();
  };

  const handleAcceptFriendship = () => {
    if (!friendshipState?.friendshipId) return;
    replyFriendship({
      friendshipId: friendshipState.friendshipId,
      accept: true,
    })
  };

  const handleRejectFriendship = () => {
    if (!friendshipState?.friendshipId) return;
    replyFriendship({
      friendshipId: friendshipState.friendshipId,
      accept: false,
    })
  };

  const handleGoToRoute = (route: string) => {
    navigate(route);
  };

  const itemsList = [
    {
      text: "recados", avatar: "📝", onClick: () => handleGoToRoute(`/profile/${userId}/scraps`)
    },
    { text: "álbum", avatar: "📷", onClick: () => handleGoToRoute("albums") },
    { text: "vídeos", avatar: "📹", onClick: () => handleGoToRoute("videos") },
    {
      text: "escrever depoimento",
      avatar: "🖊️",
      onClick: () => handleGoToRoute("testimonials"),
    },
    { text: "reportar", avatar: "🚫", onClick: () => null },
  ];

  if (friendshipState) {
    if (friendshipState.amIFriend) {
      itemsList.unshift({ text: "desfazer amizade", avatar: "➖", onClick: () => handleRemoveFriend() })
    } else if (friendshipState.isPendingMyReply) {
      itemsList.unshift({ text: "pedido enviado", avatar: "🕛", onClick: () => handleDoNothing() })
    } else if (friendshipState.isPendingTheirReply) {
      itemsList.unshift({ text: "aceitar amizade", avatar: "✅", onClick: () => handleAcceptFriendship() })
      itemsList.unshift({ text: "recusar amizade", avatar: "❌", onClick: () => handleRejectFriendship() })
    } else {
      itemsList.unshift({ text: "adicionar", avatar: "➕", onClick: () => handleAddFriend() })
    }
  }

  return (
    <Box>
      <LateralButtonList items={itemsList} />
    </Box>
  );
};
