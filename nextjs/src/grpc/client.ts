import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import path from 'path';
import { ProtoGrpcType } from './rpc/chat';

const packageDefinition = protoLoader.loadSync(path.resolve(process.cwd(), 'proto', 'chat.proto'));

const proto = grpc.loadPackageDefinition(packageDefinition) as unknown as ProtoGrpcType;

export const chatClient = new proto.pb.ChatService('localhost:50051', grpc.credentials.createInsecure());
