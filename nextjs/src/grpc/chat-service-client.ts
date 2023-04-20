import { Metadata } from '@grpc/grpc-js';
import { chatClient } from './client';
import { ChatServiceClient as GrpcChatServiceClient } from './rpc/pb/ChatService';

type ChatStreamData = {
	chat_id?: string;
	user_id: string;
	message: string;
};

export class ChatServiceClient {
	private token = '123456';
	constructor(private grpcClient: GrpcChatServiceClient) {}

	public chatStream(data: ChatStreamData) {
		const metadata = new Metadata();
		metadata.set('authorization', this.token);

		const stream = this.grpcClient.chatStream(
			{
				chatId: data.chat_id,
				userId: data.user_id,
				userMessage: data.message
			},
			metadata
		);

		stream.on('data', (data) => {
			console.log(data);
		});

		stream.on('error', (error) => {
			console.log(error);
		});

		stream.on('end', () => {
			console.log('end');
		});

		return stream;
	}
}

export class ChatServiceClientFactory {
	public static create() {
		return new ChatServiceClient(chatClient);
	}
}
