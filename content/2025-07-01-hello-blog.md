# Ditching GoMobile: Building Our Own Native Bridge

For the longest time, we used `gomobile` to bridge our Go-based core with Android and iOS. It seemed like the perfect solution—generate bindings, avoid boilerplate, and focus on logic. But as our product matured, `gomobile` started showing its cracks.

## Why We Used GoMobile in the First Place

- Automatic generation of Java and Objective-C bindings
- Easy method export from Go
- Quick prototype integration for cross-platform logic

But then reality hit.

## The Problems

- **Missing modern language support:** No Kotlin coroutines, no Swift async/await
- **Limited threading model:** Everything ran on a single locked thread
- **Opaque crashes:** Stack traces from `gomobile` were cryptic and unhelpful
- **No support for streaming, async calls, or cancellation**
- **Hard to debug JNI/ObjC glue code**

Worst of all: **gomobile was abandoned**. No meaningful updates in years, yet we were depending on it.

## What We Built Instead

We now use a **custom native bridge**, powered by:

- ✨ **Protobuf for typed request/response definitions**
- 🧵 **Cancellable async calls** (Kotlin coroutines, Swift async/await)
- 🔄 **Bidirectional communication** (platform ↔ core)
- 🧩 **Manually managed JNI and Objective-C glue code**
- 🔍 **Full control over threading and memory**

The bridge consists of:

1. A **Go core** that registers all callable methods
2. A **native dispatcher** on each platform that routes Protobuf messages
3. Typed messages and responses, fully testable
4. Direct hooks into Go’s panic handler (with stacktrace propagation)

## The Benefits

- 🔥 No more silent crashes
- 📱 Proper threading and cancellation
- 📦 Much smaller binary size
- 🛠 Full debugging capabilities
- 🚀 Faster startup and method invocation

## Example: Calling from Kotlin

```kotlin
val request = GetMediaRequest.newBuilder().build()
val response = nativeBridge.callGo(METHOD_GET_MEDIA, request.toByteArray())
```
