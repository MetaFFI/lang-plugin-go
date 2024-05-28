#ifndef TRAVERSE_CONSTRUCT_GO_CALLBACKS_H
#define TRAVERSE_CONSTRUCT_GO_CALLBACKS_H

void onFloat64(const metaffi_size* index, metaffi_size indexSize, metaffi_float64 val, void* context);
void onFloat32(const metaffi_size* index, metaffi_size indexSize, metaffi_float32 val, void* context);
void onInt8(const metaffi_size* index, metaffi_size indexSize, metaffi_int8 val, void* context);
void onUInt8(const metaffi_size* index, metaffi_size indexSize, metaffi_uint8 val, void* context);
void onInt16(const metaffi_size* index, metaffi_size indexSize, metaffi_int16 val, void* context);
void onUInt16(const metaffi_size* index, metaffi_size indexSize, metaffi_uint16 val, void* context);
void onInt32(const metaffi_size* index, metaffi_size indexSize, metaffi_int32 val, void* context);
void onUInt32(const metaffi_size* index, metaffi_size indexSize, metaffi_uint32 val, void* context);
void onInt64(const metaffi_size* index, metaffi_size indexSize, metaffi_int64 val, void* context);
void onUInt64(const metaffi_size* index, metaffi_size indexSize, metaffi_uint64 val, void* context);
void onBool(const metaffi_size* index, metaffi_size indexSize, metaffi_bool val, void* context);
void onChar8(const metaffi_size* index, metaffi_size indexSize, struct metaffi_char8 val, void* context);
void onString8(const metaffi_size* index, metaffi_size indexSize, metaffi_string8 val, void* context);
void onChar16(const metaffi_size* index, metaffi_size indexSize, struct metaffi_char16 val, void* context);
void onChar32(const metaffi_size* index, metaffi_size indexSize, struct metaffi_char32 val, void* context);
void onString16(const metaffi_size* index, metaffi_size indexSize, char16_t* val, void* context);
void onString32(const metaffi_size* index, metaffi_size indexSize, metaffi_string32 val, void* context);
void onHandle(const metaffi_size* index, metaffi_size indexSize, const struct cdt_metaffi_handle* val, void* context);
void onCallable(const metaffi_size* index, metaffi_size indexSize, const struct cdt_metaffi_callable* val, void* context);
void onNull(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_bool onArray(const metaffi_size* index, metaffi_size indexSize, const struct cdts* val, metaffi_int64 fixedDimensions, metaffi_type commonType, void* context);
void getRootElementsCount(void* context);
struct metaffi_type_info getTypeInfo(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_float64 getFloat64(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_float32 getFloat32(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_int8 getInt8(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_uint8 getUInt8(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_int16 getInt16(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_uint16 getUInt16(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_int32 getInt32(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_uint32 getUInt32(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_int64 getInt64(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_uint64 getUInt64(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_bool getBool(const metaffi_size* index, metaffi_size indexSize, void* context);
struct metaffi_char8 getChar8(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_string8 getString8(const metaffi_size* index, metaffi_size indexSize, metaffi_bool* free_required, void* context);
struct metaffi_char16 getChar16(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_string16 getString16(const metaffi_size* index, metaffi_size indexSize, metaffi_bool* free_required, void* context);
struct metaffi_char32 getChar32(const metaffi_size* index, metaffi_size indexSize, void* context);
metaffi_string32 getString32(const metaffi_size* index, metaffi_size indexSize, metaffi_bool* free_required, void* context);
struct cdt_metaffi_handle* getHandle(const metaffi_size* index, metaffi_size indexSize, metaffi_bool* free_required, void* context);
struct cdt_metaffi_callable* getCallable(const metaffi_size* index, metaffi_size indexSize, metaffi_bool* free_required, void* context);
metaffi_size getArrayMetadata(const metaffi_size* index, metaffi_size indexSize, metaffi_bool* isFixedDimension, metaffi_bool* is1DArray, metaffi_type* commonType, metaffi_bool* isManuallyConstructArray, void* context);
void constructCDTArray(const metaffi_size* index, metaffi_size indexSize, struct cdts* manuallyFillArray, void* context);

#endif //TRAVERSE_CONSTRUCT_GO_CALLBACKS_H