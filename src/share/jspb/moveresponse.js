/**
 * @fileoverview
 * @enhanceable
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

goog.provide('proto.pb.MoveResponse');

goog.require('jspb.Message');
goog.require('jspb.BinaryReader');
goog.require('jspb.BinaryWriter');
goog.require('proto.pb.Point');


/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.pb.MoveResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.pb.MoveResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.pb.MoveResponse.displayName = 'proto.pb.MoveResponse';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.pb.MoveResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.pb.MoveResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.pb.MoveResponse} msg The msg instance to transform.
 * @return {!Object}
 */
proto.pb.MoveResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    mapid: jspb.Message.getFieldWithDefault(msg, 1, 0),
    srcposition: (f = msg.getSrcposition()) && proto.pb.Point.toObject(includeInstance, f),
    dstposition: (f = msg.getDstposition()) && proto.pb.Point.toObject(includeInstance, f),
    userid: jspb.Message.getFieldWithDefault(msg, 4, 0),
    name: jspb.Message.getFieldWithDefault(msg, 5, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.pb.MoveResponse}
 */
proto.pb.MoveResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.pb.MoveResponse;
  return proto.pb.MoveResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.pb.MoveResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.pb.MoveResponse}
 */
proto.pb.MoveResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setMapid(value);
      break;
    case 2:
      var value = new proto.pb.Point;
      reader.readMessage(value,proto.pb.Point.deserializeBinaryFromReader);
      msg.setSrcposition(value);
      break;
    case 3:
      var value = new proto.pb.Point;
      reader.readMessage(value,proto.pb.Point.deserializeBinaryFromReader);
      msg.setDstposition(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUserid(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.pb.MoveResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.pb.MoveResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.pb.MoveResponse} message
 * @param {!jspb.BinaryWriter} writer
 */
proto.pb.MoveResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMapid();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getSrcposition();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.pb.Point.serializeBinaryToWriter
    );
  }
  f = message.getDstposition();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.pb.Point.serializeBinaryToWriter
    );
  }
  f = message.getUserid();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
};


/**
 * optional int64 mapId = 1;
 * @return {number}
 */
proto.pb.MoveResponse.prototype.getMapid = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.pb.MoveResponse.prototype.setMapid = function(value) {
  jspb.Message.setField(this, 1, value);
};


/**
 * optional Point srcPosition = 2;
 * @return {?proto.pb.Point}
 */
proto.pb.MoveResponse.prototype.getSrcposition = function() {
  return /** @type{?proto.pb.Point} */ (
    jspb.Message.getWrapperField(this, proto.pb.Point, 2));
};


/** @param {?proto.pb.Point|undefined} value */
proto.pb.MoveResponse.prototype.setSrcposition = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


proto.pb.MoveResponse.prototype.clearSrcposition = function() {
  this.setSrcposition(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.pb.MoveResponse.prototype.hasSrcposition = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Point dstPosition = 3;
 * @return {?proto.pb.Point}
 */
proto.pb.MoveResponse.prototype.getDstposition = function() {
  return /** @type{?proto.pb.Point} */ (
    jspb.Message.getWrapperField(this, proto.pb.Point, 3));
};


/** @param {?proto.pb.Point|undefined} value */
proto.pb.MoveResponse.prototype.setDstposition = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.pb.MoveResponse.prototype.clearDstposition = function() {
  this.setDstposition(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.pb.MoveResponse.prototype.hasDstposition = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional int64 userId = 4;
 * @return {number}
 */
proto.pb.MoveResponse.prototype.getUserid = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/** @param {number} value */
proto.pb.MoveResponse.prototype.setUserid = function(value) {
  jspb.Message.setField(this, 4, value);
};


/**
 * optional string name = 5;
 * @return {string}
 */
proto.pb.MoveResponse.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.pb.MoveResponse.prototype.setName = function(value) {
  jspb.Message.setField(this, 5, value);
};


