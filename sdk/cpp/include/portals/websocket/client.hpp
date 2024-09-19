#pragma once

#include <libwebsockets.h>

#include <memory>

#include "../../network/node.hpp"
#include "portal.hpp"

namespace directmq::portal::websocket {
class WebsocketClient;

class WebsocketClientCreationResult {
   public:
    std::shared_ptr<WebsocketClient> client;
    const char* error;
};

class WebsocketClient {
   private:
    lws_context_creation_info info;
    lws_protocols protocols[2] = {
        {
            "directmq.v1",
            WebsocketClient::dmqV1ProtocolHandler,
            sizeof(DmqProtocolConnectionContext),
            0,
        },
        {NULL, NULL, 0, 0} /* terminator */
    };
    DmqProtocolConnectionContext initialDmqProtoCtx;

    lws_context* context;
    lws* wsi;

    static int dmqV1ProtocolHandler(struct lws* wsi,
                                    enum lws_callback_reasons reason,
                                    void* user, void* in, size_t len) {
        DmqProtocolConnectionContext* ctx =
            static_cast<DmqProtocolConnectionContext*>(user);

        switch (reason) {
            case LWS_CALLBACK_CLIENT_ESTABLISHED: {
                lwsl_user("Client connected\n");

                user = ctx = new DmqProtocolConnectionContext(*ctx);

                break;
            }

            case LWS_CALLBACK_CLIENT_WRITEABLE: {
                lwsl_user("Client writeable\n");

                ctx->portal =
                    std::make_shared<WebsocketPortal>(wsi, ctx->dataFormat);
                ctx->edge = ctx->node->addConnectingEdge(ctx->portal);

                break;
            }

            case LWS_CALLBACK_CLIENT_RECEIVE: {
                lwsl_user("Received data: %s\n", (char*)in);

                ctx->edge->processIncomingPacket(
                    portal::Packet::fromData((uint8_t*)in, len));

                break;
            }

            case LWS_CALLBACK_CLIENT_CLOSED: {
                printf("Connection closed\n");

                ctx->node->removeEdge(ctx->portal, "websocket client closed");
                delete (DmqProtocolConnectionContext*)user;

                break;
            }

            default:
                break;
        }
        return 0;
    }

    WebsocketClient() {}

   public:
    static WebsocketClientCreationResult connect(const char* address,
                                                 const int port,
                                                 const char* path,
                                                 WsDataFormat dataFormat,
                                                 network::NetworkNode* node) {
        WebsocketClient client;

        // init lws context creation info
        memset(&client.info, 0, sizeof(info));
        client.info.port = port;
        client.info.iface = address;
        client.info.protocols = client.protocols;
        client.info.user = &client.initialDmqProtoCtx;

        // configure initial dmq proto context
        client.initialDmqProtoCtx.node = node;
        client.initialDmqProtoCtx.edge = nullptr;
        client.initialDmqProtoCtx.portal = nullptr;
        client.initialDmqProtoCtx.dataFormat = dataFormat;

        client.context = lws_create_context(&client.info);
        if (!client.context) {
            return WebsocketClientCreationResult{nullptr, "lws init failed"};
        }

        lws_client_connect_info ccinfo;
        memset(&ccinfo, 0, sizeof(ccinfo));
        ccinfo.context = client.context;
        ccinfo.address = address;
        ccinfo.port = port;
        ccinfo.path = path;
        ccinfo.host = lws_canonical_hostname(client.context);
        ccinfo.origin = "origin";
        ccinfo.protocol = client.protocols[0].name;
        ccinfo.ssl_connection = 0;

        client.wsi = lws_client_connect_via_info(&ccinfo);
        if (!client.wsi) {
            lws_context_destroy(client.context);
            return WebsocketClientCreationResult{nullptr,
                                                 "lws client connect failed"};
        }
    }

    ~WebsocketClient() { lws_context_destroy(context); }

    void run(int timeoutMs = 0) { lws_service(context, timeoutMs); }
};
}  // namespace directmq::portal::websocket
