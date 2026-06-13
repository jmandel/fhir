import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.List;
/**
 * The committed bootstrap for the tx.lock-driven spec build - the gradle-wrapper pattern:
 * this tiny jar is stable and lives in the repo; the real tooling is named (url + sha256) by
 * eng/future/kindling-wrapper.properties and fetched/verified/cached automatically. The whole developer experience is:
 *
 *   java -jar eng/future/launch.jar build .
 *
 * No bash, no curl, no python; identical on Windows/macOS/Linux. Dependency footprint: a JDK.
 * The launcher downloads the pinned tooling jar into ~/.fhir/tools/<sha256>.jar (verifying the
 * hash, refusing anything else), then runs it in a child JVM with a sensible default heap
 * (-Xmx12g, override with the HEAP env var, e.g. HEAP=11g) and passes your arguments through
 * to org.hl7.fhir.tools.publisher.SpecBuild.
 */
public class Launcher {

  public static void main(String[] args) throws Exception {
    File props = new File("eng/future/kindling-wrapper.properties");
    if (!props.exists()) {
      System.err.println("no eng/future/kindling-wrapper.properties - run from the spec checkout root");
      System.exit(2);
    }
    java.util.Properties pin = new java.util.Properties();
    try (java.io.FileInputStream in = new java.io.FileInputStream(props)) {
      pin.load(in);
    }
    String url = pin.getProperty("toolUrl");
    String sha = pin.getProperty("toolSha256");
    if (url == null || sha == null || !sha.matches("[0-9a-f]{64}")) {
      System.err.println("kindling-wrapper.properties has no usable toolUrl / toolSha256");
      System.exit(2);
    }

    File store = Paths.get(System.getProperty("user.home"), ".fhir", "tools").toFile();
    store.mkdirs();
    File jar = new File(store, sha + ".jar");
    if (!jar.exists() || !sha.equals(sha256(jar))) {
      System.out.println("fetching build tooling (" + url + ") ...");
      File part = new File(store, sha + ".jar.download");
      download(url, part);
      if (!sha.equals(sha256(part))) {
        part.delete();
        System.err.println("downloaded tooling does not match tx.lock tooling.sha256 - refusing to run it");
        System.exit(1);
      }
      if (jar.exists()) {
        jar.delete();
      }
      if (!part.renameTo(jar)) {
        System.err.println("could not install tooling into " + jar);
        System.exit(1);
      }
    }

    List<String> cmd = new ArrayList<>();
    cmd.add(Paths.get(System.getProperty("java.home"), "bin",
        System.getProperty("os.name").toLowerCase().contains("win") ? "java.exe" : "java").toString());
    cmd.add("-Xmx" + System.getenv().getOrDefault("HEAP", "12g"));
    String opts = System.getenv("JAVA_OPTS");
    if (opts != null) {
      for (String o : opts.trim().split("\\s+")) {
        cmd.add(o);
      }
    }
    cmd.add("-cp");
    cmd.add(jar.getAbsolutePath());
    cmd.add("org.hl7.fhir.tools.publisher.SpecBuild");
    for (String a : args) {
      cmd.add(a);
    }
    Process p = new ProcessBuilder(cmd).inheritIO().start();
    System.exit(p.waitFor());
  }

  private static void download(String url, File dest) throws IOException {
    HttpURLConnection.setFollowRedirects(true);
    HttpURLConnection c = (HttpURLConnection) new URL(url).openConnection();
    c.setConnectTimeout(15000);
    c.setReadTimeout(60000);
    int code = c.getResponseCode();
    // follow one cross-protocol/cross-host redirect chain manually if the JDK stopped early -
    // Location headers are used AS-IS (never decoded: signed URLs carry encoded parameters)
    int hops = 0;
    while (code >= 301 && code <= 308 && hops++ < 5) {
      String loc = c.getHeaderField("Location");
      if (loc == null) {
        throw new IOException("redirect without Location from " + url);
      }
      c = (HttpURLConnection) new URL(new URL(url), loc).openConnection();
      c.setConnectTimeout(15000);
      c.setReadTimeout(60000);
      code = c.getResponseCode();
    }
    if (code != 200) {
      throw new IOException("HTTP " + code + " fetching " + url);
    }
    try (InputStream in = c.getInputStream(); FileOutputStream out = new FileOutputStream(dest)) {
      byte[] buf = new byte[65536];
      int n;
      while ((n = in.read(buf)) > 0) {
        out.write(buf, 0, n);
      }
    }
  }

  private static String sha256(File f) throws Exception {
    MessageDigest md = MessageDigest.getInstance("SHA-256");
    byte[] hash = md.digest(Files.readAllBytes(f.toPath()));
    StringBuilder b = new StringBuilder();
    for (byte x : hash) {
      b.append(String.format("%02x", x));
    }
    return b.toString();
  }
}
