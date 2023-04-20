import { prisma } from '@/app/prisma/prisma';
import { NextResponse } from 'next/server';

type Params = {
	params: {
		chatId: string;
	};
};

export async function GET(req: Request, { params }: Params) {
	const { chatId } = params;

	const messages = await prisma.message.findMany({
		where: {
			chat_id: chatId
		},
		orderBy: {
			created_at: 'desc'
		}
	});

	return NextResponse.json(messages);
}

export async function POST(req: Request, { params }: Params) {
	const { chatId } = params;

	const chat = await prisma.chat.findUnique({
		where: {
			id: chatId
		}
	});

	if (!chat) {
		return new Response('Chat not found', { status: 404 });
	}

	const body = await req.json();

	const messageCreated = await prisma.message.create({
		data: {
			content: body.message,
			chat_id: chatId
		}
	});

	return NextResponse.json(messageCreated);
}
