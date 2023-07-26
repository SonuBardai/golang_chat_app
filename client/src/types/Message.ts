type Levels = "info" | "error" | null;
export interface RenderableMessage {
  message: string;
  user: string;
  level: Levels;
}

type MessageType = "newUser" | "broadcastMessage";

interface BaseMessage {
  messageType: MessageType;
}

export interface BroadcastMessage extends BaseMessage {
  username: string;
  content: string;
}

export interface NewUserMessage extends BaseMessage {
  username: string;
}

export type ServerMessage = BroadcastMessage | NewUserMessage;
