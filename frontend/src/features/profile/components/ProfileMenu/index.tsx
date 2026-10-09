import React from "react";
import { LateralButtonList } from "components/LateralButtonList";

import { Box } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useMe, useMeAmFriend } from "api/me/queries";
import { useAddFriend, useRemoveFriend, useReplyFriendship } from "api/friendship/queries";

interface IProfileMenuProps {
  userId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ userId }) => {
  const navigate = useNavigate();
  const { data: me } = useMe();
  const { data: friendshipState } = useMeAmFriend(userId);
  const { mutate: replyFriendship, isPending: isPendingReply } = useReplyFriendship();
  const { mutate: removeFriend, isPending: isPendingRemoval } = useRemoveFriend(userId);
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
      text: "perfil", avatar: "🙂", onClick: () => handleGoToRoute(`/profile/${userId}`)
    },
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
    { text: "reportar", avatar: "🚫" },
  ];

  if (friendshipState) {
    if (friendshipState.amIFriend) {
      itemsList.unshift({ text: "desfazer amizade", avatar: "➖", onClick: () => handleRemoveFriend() })
    } else if (friendshipState.isPendingTheirReply) {
      itemsList.unshift({ text: "pedido enviado", avatar: "🕛", onClick: () => handleDoNothing() })
    } else if (friendshipState.isPendingMyReply) {
      itemsList.unshift({ text: "aceitar amizade", avatar: "✅", onClick: () => handleAcceptFriendship() })
      itemsList.unshift({ text: "recusar amizade", avatar: "❌", onClick: () => handleRejectFriendship() })
    } else {
      itemsList.unshift({ text: "adicionar", avatar: "➕", onClick: () => handleAddFriend() })
    }
  }

  // Last unshift, so it's always the first option on your own profile.
  if (me?.id === userId) {
    itemsList.unshift({ text: "editar perfil", avatar: "✏️", onClick: () => handleGoToRoute(`/profile/${userId}/edit`) })
  }

  return (
    <Box>
      <LateralButtonList items={itemsList} />
    </Box>
  );
};
