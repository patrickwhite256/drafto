/**
 * Code generated by protoc-gen-twirp_js v5.0.0, DO NOT EDIT.
 * source: service.proto
 */
// import our twirp js library dependency
var createClient = require("twirp");
// import our protobuf definitions
var pb = require("./service_pb.js");
Object.assign(module.exports, pb);

/**
 * Creates a new DraftoClient
 */
module.exports.createDraftoClient = function(baseurl, extraHeaders, useJSON) {
    var rpc = createClient(baseurl, "patrickwhite256.drafto.Drafto", "v5.0.0",  useJSON, extraHeaders === undefined ? {} : extraHeaders);
    return {
        newDraft: function(data) { return rpc("NewDraft", data, pb.NewDraftResp); },
        getSeat: function(data) { return rpc("GetSeat", data, pb.GetSeatResp); },
        makeSelection: function(data) { return rpc("MakeSelection", data, pb.MakeSelectionResp); },
        getDraftStatus: function(data) { return rpc("GetDraftStatus", data, pb.GetDraftStatusResp); },
        takeSeat: function(data) { return rpc("TakeSeat", data, pb.TakeSeatResp); },
        getCurrentUser: function(data) { return rpc("GetCurrentUser", data, pb.GetCurrentUserResp); }
    }
}

