#pragma once

#include <libwebsockets.h>

#include <vector>

#include "../../network/edge/network_edge.hpp"
#include "../../network/node.hpp"
#include "../../portal.hpp"

namespace directmq::portal::websocket {
enum class WsDataFormat : int { BINARY = 1, TEXT = 0 };

struct DmqProtocolConnectionContext {
    network::NetworkNode* node;
    std::shared_ptr<WebsocketPortal> portal;
    std::shared_ptr<network::edge::NetworkEdge> edge;
    WsDataFormat dataFormat;
};

class WebsocketWriter : public DataWriter {
   private:
    lws* wsi;

    WsDataFormat dataFormat;

    const size_t messageSize;
    size_t writtenBytes;

   public:
    WebsocketWriter(lws* wsi, WsDataFormat dataFormat, size_t messageSize)
        : wsi(wsi), dataFormat(dataFormat), messageSize(messageSize) {}

    bool write(const uint8_t* block, const size_t blockSize) override {
        if (writtenBytes + blockSize > messageSize) {
            throw new std::runtime_error(
                "Trying to write more data than expected");
        }

        int flags = dataFormat == WsDataFormat::BINARY ? LWS_WRITE_BINARY
                                                       : LWS_WRITE_TEXT;

        bool isFirstChunk = writtenBytes == 0;
        if (!isFirstChunk) {
            flags = LWS_WRITE_CONTINUATION;
        }

        bool isLastChunk = writtenBytes + blockSize == messageSize;
        if (isLastChunk) {
            flags &= ~LWS_WRITE_NO_FIN;
        } else {
            flags |= LWS_WRITE_NO_FIN;
        }

        std::vector<uint8_t> buffer(LWS_PRE + blockSize);
        std::copy(block, block + blockSize, buffer.begin() + LWS_PRE);

        int written = lws_write(wsi, buffer.data() + LWS_PRE, blockSize,
                                (lws_write_protocol)flags);
        if (written != blockSize) {
            // TODO: instead split block in MTU size chunks and send them
            // using loop to handle truncated writes
            return false;
        }

        writtenBytes += blockSize;
        return true;
    }

    void end() override {
        // nothing to do here, as all flags
        // were sent by write method already
    }
};

class WebsocketPortal : public Portal {
   private:
    lws* wsi;
    WsDataFormat dataFormat;

   public:
    WebsocketPortal(struct lws* wsi, WsDataFormat dataFormat)
        : wsi(wsi), dataFormat(dataFormat) {}

    /**
     * Portal interface implementation
     */

    std::shared_ptr<DataWriter> beginWrite(const size_t packetSize) override {
        return std::make_shared<WebsocketWriter>(wsi, dataFormat, packetSize);
    }

    void close() override {
        lws_close_reason(wsi, LWS_CLOSE_STATUS_NORMAL, nullptr, 0);
    }
};
}  // namespace directmq::portal::websocket
