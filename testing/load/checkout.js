import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['../..'], './checkout/api/checkout/v1/checkout.proto');

export const options = {
    vus: 5,
    duration: '30m'
};

export default () => {
    client.connect('localhost:30000', {
        plaintext: true
    });

    const data = {
        "user": 1,
        "sku": 1148162,
        "count": 1
    }
    const list = {
        "user": 1,
    }

    data.count++
    const responseAdd = client.invoke('checkout.Checkout/AddToCart', data);
    check(responseAdd, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    if (responseAdd.status !== grpc.StatusOK) {
        console.log(JSON.stringify(responseAdd.error));
    }

    data.count--
    const responseDel = client.invoke('checkout.Checkout/DeleteFromCart', data);
    check(responseDel, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    if (responseDel.status !== grpc.StatusOK) {
        console.log(JSON.stringify(responseDel.error));
    }


    const responseList = client.invoke('checkout.Checkout/ListCart', list);
    check(responseList, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    if (responseList.status !== grpc.StatusOK) {
        console.log(JSON.stringify(responseList.error));
    }

    const responseDelete = client.invoke('checkout.Checkout/Purchase', list);
    check(responseDelete, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    if (responseDelete.status !== grpc.StatusOK) {
        console.log(JSON.stringify(responseDelete.error));
    }

    client.close();
    sleep(1);
};
