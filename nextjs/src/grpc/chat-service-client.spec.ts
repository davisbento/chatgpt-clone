import { ChatServiceClientFactory } from './chat-service-client';

describe('ChatServiceClient', () => {
	it('should grpc client work', (done) => {
		const client = ChatServiceClientFactory.create();

		const stream = client.chatStream({
			user_id: '1',
			message: 'Hello World'
		});

		stream.on('end', () => {
			done();
		});
	});
});
