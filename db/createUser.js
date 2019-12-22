const user = {
    user: 'root',
    pwd: 'example',
    // roles: [{role: "userAdminAnyDatabase", db: "admin"}],
    roles: [{
        role: 'readWrite',
        db: 'blog'

    }]
};

db.createUser(user);

const user2 = {
    user: 'adminroot',
    pwd: 'adminexample',
    roles: [{role: "userAdminAnyDatabase", db: "admin"}],
    // roles: [{
    //     role: 'readWrite',
    //     db: 'blog'
    //
    // }]
};

db.createUser(user2);
