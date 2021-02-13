desc: updates auth-config of eks cluster every minute. it reads aws iam groups and mapped them to eks roles. it is forked from mentioned repos and aws config structure changed to work with oidc.

example:
Assuming you have 2 IAM groups devops and devs, and you want to give devops's members system:masters role and devs a developer role, you can do that via:

./app --iam-k8s-group=devops::system:masters,devs::developer
./app --iam-k8s-group=devops::system:masters,devs::developer|manager (supports multiple kubernetes roles)

Thanks to:
    https://github.com/MindTickle/iam-eks-user-mapper
    https://github.com/ygrene/iam-eks-user-mapper