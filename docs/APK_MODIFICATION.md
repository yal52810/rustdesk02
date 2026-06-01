# RustDesk APK 修改指南 — 国产手机防封

## 问题背景

小米（MIUI）、OPPO（ColorOS）、Vivo（OriginOS）、华为（HarmonyOS）等国产手机系统会对远程桌面类应用进行主动拦截：

- **包名黑名单**：系统内置黑名单，已知远程桌面包名（如 `com.carriez.flutter_hbb`）会被直接标记为风险应用
- **安装拦截**：侧载 APK 时弹出"存在风险"警告，甚至直接禁止安装
- **后台限制**：即使安装成功，系统也会严格限制后台运行和通知权限

## 前提条件

RustDesk 客户端源码 **不在本仓库中**。本仓库仅包含服务端代码。

客户端源码仓库：<https://github.com/rustdesk/rustdesk>

修改 APK 需要先克隆客户端仓库：

```bash
git clone https://github.com/rustdesk/rustdesk.git
cd rustdesk
```

## 1. 修改包名（Package Name）

### 1.1 修改 build.gradle

编辑 `android/app/build.gradle`，修改 `applicationId`：

```groovy
android {
    defaultConfig {
        // 原值: applicationId "com.carriez.flutter_hbb"
        applicationId "com.yourcompany.remoteassist"  // 改为你的包名
        // ...
    }
}
```

### 1.2 修改 AndroidManifest.xml

编辑 `android/app/src/main/AndroidManifest.xml`，将 `package` 属性改为新包名：

```xml
<!-- 原值: package="com.carriez.flutter_hbb" -->
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="com.yourcompany.remoteassist">
```

### 1.3 重命名 Kotlin 目录

```bash
# 原目录结构
android/app/src/main/kotlin/com/carriez/flutter_hbb/

# 改为新的目录结构
android/app/src/main/kotlin/com/yourcompany/remoteassist/
```

```bash
# Linux/Mac
cd android/app/src/main/kotlin
mkdir -p com/yourcompany/remoteassist
cp -r com/carriez/flutter_hbb/* com/yourcompany/remoteassist/
rm -rf com/carriez

# Windows PowerShell
cd android\app\src\main\kotlin
New-Item -ItemType Directory -Force -Path com\yourcompany\remoteassist
Copy-Item -Recurse com\carriez\flutter_hbb\* com\yourcompany\remoteassist\
Remove-Item -Recurse com\carriez
```

### 1.4 修改所有 Kotlin 文件中的 package 声明

```bash
# 批量替换所有 .kt 文件中的包名
grep -rl "com.carriez.flutter_hbb" android/app/src/main/kotlin/ | \
  xargs sed -i 's/com\.carriez\.flutter_hbb/com.yourcompany.remoteassist/g'
```

### 1.5 修改 Flutter 侧的包名引用

检查并修改以下文件中可能存在的包名引用：

- `lib/main.dart`
- `pubspec.yaml`（修改 `name` 字段）

## 2. 修改应用名称和图标

### 2.1 修改应用名称

编辑 `android/app/src/main/res/values/strings.xml`：

```xml
<resources>
    <!-- 原值: <string name="app_name">RustDesk</string> -->
    <string name="app_name">远程助手</string>  <!-- 改为一个不引人注意的名字 -->
</resources>
```

### 2.2 替换应用图标

替换以下目录中的 `ic_launcher.png`（及其他图标文件）：

- `android/app/src/main/res/mipmap-mdpi/`
- `android/app/src/main/res/mipmap-hdpi/`
- `android/app/src/main/res/mipmap-xhdpi/`
- `android/app/src/main/res/mipmap-xxhdpi/`
- `android/app/src/main/res/mipmap-xxxhdpi/`

推荐使用自适应图标（Adaptive Icon）格式。

### 2.3 修改通知图标

替换 `android/app/src/main/res/drawable/` 中的通知相关图标。

## 3. 生成签名密钥并配置

### 3.1 生成密钥

```bash
keytool -genkey -v \
  -keystore your-key.keystore \
  -alias your-alias \
  -keyalg RSA \
  -keysize 2048 \
  -validity 10000
```

记住你设置的密钥库密码和密钥密码。

### 3.2 创建 key.properties

在 `android/` 目录下创建 `key.properties`：

```properties
storePassword=你的密钥库密码
keyPassword=你的密钥密码
keyAlias=your-alias
storeFile=../your-key.keystore
```

### 3.3 配置签名

编辑 `android/app/build.gradle`，在 `android` 块之前添加：

```groovy
def keystoreProperties = new Properties()
def keystorePropertiesFile = rootProject.file('key.properties')
if (keystorePropertiesFile.exists()) {
    keystoreProperties.load(new FileInputStream(keystorePropertiesFile))
}

android {
    // ...
    signingConfigs {
        release {
            keyAlias keystoreProperties['keyAlias']
            keyPassword keystoreProperties['keyPassword']
            storeFile keystoreProperties['storeFile'] ? file(keystoreProperties['storeFile']) : null
            storePassword keystoreProperties['storePassword']
        }
    }
    buildTypes {
        release {
            signingConfig signingConfigs.release
        }
    }
}
```

## 4. 启用 ProGuard/R8 混淆

编辑 `android/app/build.gradle`：

```groovy
android {
    buildTypes {
        release {
            minifyEnabled true
            shrinkResources true
            proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
        }
    }
}
```

创建 `android/app/proguard-rules.pro`，添加基本规则：

```proguard
# 保留 Flutter 相关
-keep class io.flutter.** { *; }
-keep class com.yourcompany.remoteassist.** { *; }

# 混淆其他类
-keepattributes Signature
-keepattributes *Annotation*
```

## 5. 构建 APK

```bash
cd rustdesk
flutter clean
flutter pub get
flutter build apk --release
```

输出 APK 位于：`build/app/outputs/flutter-apk/app-release.apk`

## 6. 各厂商特殊说明

| 厂商 | 系统 | 额外注意事项 |
|------|------|-------------|
| 小米 | MIUI | 需在"设置 → 应用设置 → 授权管理"中关闭"纯净模式"才能侧载。如果仍然被拦截，尝试使用不同的包名前缀（避免包含 `remote`、`vnc`、`desktop` 等关键词） |
| OPPO | ColorOS | 关闭"安全守护"中的"安装检查"。OPPO 对包含 `desktop`、`remote` 关键词的包名敏感 |
| Vivo | OriginOS | 关闭"安全检测"。Vivo 可能彻底禁止特定分类的 APK 安装 |
| 华为 | HarmonyOS | 关闭"纯净模式"。华为对未上架应用有严格的安装限制 |

## 7. 持续规避检测

### 7.1 定期更换包名

建议准备多个包名变体，当一个被列入黑名单后切换使用。

### 7.2 修改 User-Agent

在 Flutter 网络层，修改 HTTP 请求的 User-Agent，避免暴露 RustDesk 特征字符串。

### 7.3 裁剪不需要的权限

移除 `AndroidManifest.xml` 中不必要的权限声明，减少系统审查标记。

### 7.4 移除调试信息

确保 `release` 构建中 `android:debuggable="false"`，避免被识别为调试版本。

## 8. 验证清单

- [ ] APK 能在目标手机上成功安装（无拦截提示）
- [ ] 应用能正常启动
- [ ] 能连接到自定义服务器（ID 服务器 + 中继服务器）
- [ ] WSS 连接正常（通过 443 端口绕过公司防火墙）
- [ ] 远程控制功能正常
- [ ] 文件传输功能正常
- [ ] 后台运行不被系统杀死
- [ ] 通知权限正常
