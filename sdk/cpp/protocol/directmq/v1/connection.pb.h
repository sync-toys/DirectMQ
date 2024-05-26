/* Automatically generated nanopb header */
/* Generated by nanopb-0.4.8 */

#ifndef PB_DIRECTMQ_V1_DIRECTMQ_V1_CONNECTION_PB_H_INCLUDED
#define PB_DIRECTMQ_V1_DIRECTMQ_V1_CONNECTION_PB_H_INCLUDED
#include <pb.h>

#if PB_PROTO_HEADER_VERSION != 40
#error Regenerate this file with the current version of nanopb generator.
#endif

/* Struct definitions */
typedef struct _directmq_v1_SupportedProtocolVersions {
    pb_callback_t supported_protocol_versions;
} directmq_v1_SupportedProtocolVersions;

typedef struct _directmq_v1_InitConnection {
    uint64_t max_message_size;
} directmq_v1_InitConnection;

typedef struct _directmq_v1_ConnectionAccepted {
    uint64_t max_message_size;
} directmq_v1_ConnectionAccepted;

typedef struct _directmq_v1_GracefullyClose {
    pb_callback_t reason;
} directmq_v1_GracefullyClose;

typedef struct _directmq_v1_TerminateNetwork {
    pb_callback_t reason;
} directmq_v1_TerminateNetwork;


#ifdef __cplusplus
extern "C" {
#endif

/* Initializer values for message structs */
#define directmq_v1_SupportedProtocolVersions_init_default {{{NULL}, NULL}}
#define directmq_v1_InitConnection_init_default  {0}
#define directmq_v1_ConnectionAccepted_init_default {0}
#define directmq_v1_GracefullyClose_init_default {{{NULL}, NULL}}
#define directmq_v1_TerminateNetwork_init_default {{{NULL}, NULL}}
#define directmq_v1_SupportedProtocolVersions_init_zero {{{NULL}, NULL}}
#define directmq_v1_InitConnection_init_zero     {0}
#define directmq_v1_ConnectionAccepted_init_zero {0}
#define directmq_v1_GracefullyClose_init_zero    {{{NULL}, NULL}}
#define directmq_v1_TerminateNetwork_init_zero   {{{NULL}, NULL}}

/* Field tags (for use in manual encoding/decoding) */
#define directmq_v1_SupportedProtocolVersions_supported_protocol_versions_tag 1
#define directmq_v1_InitConnection_max_message_size_tag 1
#define directmq_v1_ConnectionAccepted_max_message_size_tag 1
#define directmq_v1_GracefullyClose_reason_tag   1
#define directmq_v1_TerminateNetwork_reason_tag  1

/* Struct field encoding specification for nanopb */
#define directmq_v1_SupportedProtocolVersions_FIELDLIST(X, a) \
X(a, CALLBACK, REPEATED, UINT32,   supported_protocol_versions,   1)
#define directmq_v1_SupportedProtocolVersions_CALLBACK pb_default_field_callback
#define directmq_v1_SupportedProtocolVersions_DEFAULT NULL

#define directmq_v1_InitConnection_FIELDLIST(X, a) \
X(a, STATIC,   SINGULAR, UINT64,   max_message_size,   1)
#define directmq_v1_InitConnection_CALLBACK NULL
#define directmq_v1_InitConnection_DEFAULT NULL

#define directmq_v1_ConnectionAccepted_FIELDLIST(X, a) \
X(a, STATIC,   SINGULAR, UINT64,   max_message_size,   1)
#define directmq_v1_ConnectionAccepted_CALLBACK NULL
#define directmq_v1_ConnectionAccepted_DEFAULT NULL

#define directmq_v1_GracefullyClose_FIELDLIST(X, a) \
X(a, CALLBACK, SINGULAR, STRING,   reason,            1)
#define directmq_v1_GracefullyClose_CALLBACK pb_default_field_callback
#define directmq_v1_GracefullyClose_DEFAULT NULL

#define directmq_v1_TerminateNetwork_FIELDLIST(X, a) \
X(a, CALLBACK, SINGULAR, STRING,   reason,            1)
#define directmq_v1_TerminateNetwork_CALLBACK pb_default_field_callback
#define directmq_v1_TerminateNetwork_DEFAULT NULL

extern const pb_msgdesc_t directmq_v1_SupportedProtocolVersions_msg;
extern const pb_msgdesc_t directmq_v1_InitConnection_msg;
extern const pb_msgdesc_t directmq_v1_ConnectionAccepted_msg;
extern const pb_msgdesc_t directmq_v1_GracefullyClose_msg;
extern const pb_msgdesc_t directmq_v1_TerminateNetwork_msg;

/* Defines for backwards compatibility with code written before nanopb-0.4.0 */
#define directmq_v1_SupportedProtocolVersions_fields &directmq_v1_SupportedProtocolVersions_msg
#define directmq_v1_InitConnection_fields &directmq_v1_InitConnection_msg
#define directmq_v1_ConnectionAccepted_fields &directmq_v1_ConnectionAccepted_msg
#define directmq_v1_GracefullyClose_fields &directmq_v1_GracefullyClose_msg
#define directmq_v1_TerminateNetwork_fields &directmq_v1_TerminateNetwork_msg

/* Maximum encoded size of messages (where known) */
/* directmq_v1_SupportedProtocolVersions_size depends on runtime parameters */
/* directmq_v1_GracefullyClose_size depends on runtime parameters */
/* directmq_v1_TerminateNetwork_size depends on runtime parameters */
#define DIRECTMQ_V1_DIRECTMQ_V1_CONNECTION_PB_H_MAX_SIZE directmq_v1_InitConnection_size
#define directmq_v1_ConnectionAccepted_size      11
#define directmq_v1_InitConnection_size          11

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif
