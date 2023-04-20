import { prisma } from '@/app/prisma/prisma';
import { NextResponse } from 'next/server';

export async function POST(req: Request) {
	const chatCreated = await prisma.chat.create({
		data: {}
	});

	return NextResponse.json({ id: chatCreated.id });
}

export async function GET(req: Request) {
	const chats = await prisma.chat.findMany();

	return NextResponse.json(chats);
}
