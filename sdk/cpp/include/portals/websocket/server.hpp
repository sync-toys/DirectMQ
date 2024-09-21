#pragma once

#include <libwebsockets.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "../../network/edge/network_edge.hpp"
#include "../../network/node.hpp"
#include "../../portal.hpp"
#include "portal.hpp"
#include "runnable.hpp"

namespace directmq::portal::websocket {
class WebsocketServer;

struct WebsocketServerCreationResult {
    std::shared_ptr<WebsocketServer> instance;
    const char* error;
};

class WebsocketServer : public Runnable {
   private:
    lws_context_creation_info info;
    lws_protocols protocols[2] = {
        {
            "directmq.v1",
            WebsocketServer::dmqV1ProtocolHandler,
            sizeof(DmqProtocolConnectionContext),
            0,
        },
        {NULL, NULL, 0, 0} /* terminator */
    };
    DmqProtocolConnectionContext initialDmqProtoCtx;
    lws_context* context;

    WebsocketServer() {}

    static int dmqV1ProtocolHandler(struct lws* wsi,
                                    enum lws_callback_reasons reason,
                                    void* user, void* in, size_t len) {
        DmqProtocolConnectionContext* ctx =
            static_cast<DmqProtocolConnectionContext*>(user);

        switch (reason) {
            case LWS_CALLBACK_ESTABLISHED: {
                printf("Connection established\n");

                user = ctx = new DmqProtocolConnectionContext(*ctx);

                ctx->portal =
                    std::make_shared<WebsocketPortal>(wsi, ctx->dataFormat);
                ctx->edge = ctx->edgeManager->addListeningEdge(ctx->portal);

                break;
            }

            case LWS_CALLBACK_RECEIVE: {
                printf("Received data: %s\n", (char*)in);

                ctx->edge->processIncomingPacket(
                    portal::Packet::fromData((uint8_t*)in, len));

                break;
            }

            case LWS_CALLBACK_CLOSED: {
                printf("Connection closed\n");

                ctx->edgeManager->removeEdge(ctx->portal,
                                             "websocket connection closed");
                delete (DmqProtocolConnectionContext*)user;

                break;
            }

            default:
                break;
        }

        return 0;
    }

   public:
    static WebsocketServerCreationResult create(
        const char* address, const int port, WsDataFormat dataFormat,
        network::EdgeManager* edgeManager) {
        WebsocketServer server;

        // init lws context creation info
        memset(&server.info, 0, sizeof(info));
        server.info.port = port;
        server.info.iface = address;
        server.info.protocols = server.protocols;
        server.info.user = &server.initialDmqProtoCtx;

        // configure initial dmq proto context
        server.initialDmqProtoCtx.edgeManager = edgeManager;
        server.initialDmqProtoCtx.edge = nullptr;
        server.initialDmqProtoCtx.portal = nullptr;
        server.initialDmqProtoCtx.dataFormat = dataFormat;

        server.context = lws_create_context(&server.info);
        if (!server.context) {
            return WebsocketServerCreationResult{nullptr, "lws init failed"};
        }

        return WebsocketServerCreationResult{
            std::make_shared<WebsocketServer>(server), nullptr};
    }

    ~WebsocketServer() { lws_context_destroy(context); }

    void run(int timeoutMs = 0) override { lws_service(context, timeoutMs); }
};
}  // namespace directmq::portal::websocket
