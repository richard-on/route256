import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['../..'], './checkout/api/checkout/v1/checkout.proto');

export const options = {
    vus: 2000,
    duration: '10m'
};

export default () => {
    client.connect('localhost:30000'/*'checkout.route256.richardhere.dev:443'*/, {
        plaintext: true
    });

    const data = {
        "user": 1,
        "sku": 33165704,
        "count": 1
    }
    const list = {
        "user": 1,
    }

    const responseAdd = client.invoke('checkout.Checkout/AddToCart', data);
    check(responseAdd, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    if (responseAdd.status !== grpc.StatusOK) {
        console.log(JSON.stringify(responseAdd.error));
    }

    /*const responseList = client.invoke('checkout.Checkout/ListCart', list);
    check(responseList, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    if (responseList.status !== grpc.StatusOK) {
        console.log(JSON.stringify(responseList.error));
    }*/

    const responseDelete = client.invoke('checkout.Checkout/DeleteFromCart', data);
    check(responseDelete, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    if (responseDelete.status !== grpc.StatusOK) {
        console.log(JSON.stringify(responseDelete.error));
    }

    client.close();
    sleep(1);
};
