import { prisma } from '@/app/prisma/prisma';
import { NextResponse } from 'next/server';

export async function POST(req: Request) {
	const body = await req.json();

	const chatCreated = await prisma.chat.create({
		data: {
			messages: {
				create: {
					content: body.message
				}
			}
		},
		select: {
			id: true,
			messages: true
		}
	});

	return NextResponse.json({ id: chatCreated.id });
}

export async function GET(req: Request) {
	const chats = await prisma.chat.findMany({
		select: {
			id: true,
			messages: {
				orderBy: {
					created_at: 'asc'
				},
				take: 1
			}
		},
		orderBy: {
			created_at: 'desc'
		}
	});

	return NextResponse.json(chats);
}
